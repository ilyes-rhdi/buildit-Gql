package services

import (
	"context"
	"errors"

	"github.com/ilyes-rhdi/buildit-Gql/internal/models"
	"gorm.io/gorm"
)

type PageService struct {
	ws *WorkspaceService
}

func NewPageService() *PageService {
	return &PageService{ws: NewWorkspaceService()}
}

func (s *PageService) CreatePage(workspaceID, requesterID string, parentPageID *string, title string) (*models.Page, error) {
	ctx := context.Background()

	ok, err := s.ws.isMember(ctx, workspaceID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not authorized")
	}

	if title == "" {
		title = "Untitled"
	}

	p := &models.Page{
		WorkspaceID:  workspaceID,
		ParentPageID: parentPageID,
		Title:        title,
		CreatedByID:  requesterID,
		Archived:     false,
	}

	if err := getDB().WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PageService) GetPage(pageID, requesterID string) (*models.Page, error) {
	ctx := context.Background()

	var p models.Page
	if err := getDB().WithContext(ctx).First(&p, "id = ?", pageID).Error; err != nil {
		return nil, err
	}

	ok, err := s.ws.isMember(ctx, p.WorkspaceID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not authorized")
	}

	return &p, nil
}

func (s *PageService) ListPages(workspaceID, requesterID string, parentPageID *string) ([]models.Page, error) {
	ctx := context.Background()

	ok, err := s.ws.isMember(ctx, workspaceID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not authorized")
	}

	q := getDB().WithContext(ctx).Where("workspace_id = ?", workspaceID).Where("archived = ?", false)
	if parentPageID == nil {
		q = q.Where("parent_page_id IS NULL")
	} else {
		q = q.Where("parent_page_id = ?", *parentPageID)
	}

	var pages []models.Page
	if err := q.Order(`"created_at" desc`).Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

type UpdatePageInput struct {
	Title    *string
	Icon     *string
	Cover    *string
	Archived *bool
	ParentID **string // si tu veux allow move dans update
}

func (s *PageService) UpdatePage(pageID, requesterID string, in UpdatePageInput) (*models.Page, error) {
	ctx := context.Background()

	var p models.Page
	if err := getDB().WithContext(ctx).First(&p, "id = ?", pageID).Error; err != nil {
		return nil, err
	}

	ok, err := s.ws.isMember(ctx, p.WorkspaceID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("not authorized")
	}

	updates := map[string]any{}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Icon != nil {
		updates["icon"] = *in.Icon
	}
	if in.Cover != nil {
		updates["cover"] = *in.Cover
	}
	if in.Archived != nil {
		updates["archived"] = *in.Archived
	}
	if in.ParentID != nil {
		updates["parent_page_id"] = *in.ParentID // peut être nil
	}

	if len(updates) == 0 {
		return &p, nil
	}

	if err := getDB().WithContext(ctx).Model(&models.Page{}).
		Where("id = ?", pageID).
		Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := getDB().WithContext(ctx).First(&p, "id = ?", pageID).Error; err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PageService) ArchivePage(pageID, requesterID string) error {
	ctx := context.Background()

	p, err := s.GetPage(pageID, requesterID)
	if err != nil {
		return err
	}

	return getDB().WithContext(ctx).
		Model(&models.Page{}).
		Where("id = ?", p.ID).
		Update("archived", true).Error
}

func (s *PageService) DeletePageHard(pageID, requesterID string) error {
	ctx := context.Background()

	p, err := s.GetPage(pageID, requesterID)
	if err != nil {
		return err
	}

	// option: supprimer blocks de la page avant
	if err := getDB().WithContext(ctx).Where("page_id = ?", p.ID).Delete(&models.Block{}).Error; err != nil {
		return err
	}

	res := getDB().WithContext(ctx).Delete(&models.Page{}, "id = ?", p.ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
