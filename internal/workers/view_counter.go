package workers

import (
	"blog_platform/internal/models"
	"log"

	"gorm.io/gorm"
)

type ViewWorker struct {
	db *gorm.DB
}

func NewViewWorker(db *gorm.DB) *ViewWorker {
	return &ViewWorker{db: db}
}

func (w *ViewWorker) ProcessViews() {
	log.Println("🔄 CRON: Processing new post views...")

	type Result struct {
		PostID   string
		NewViews int
	}
	var results []Result

	err := w.db.Model(&models.PostView{}).
		Select("post_id, count(*) as new_views").
		Where("processed = ?", false).
		Group("post_id").
		Scan(&results).Error

	if err != nil {
		log.Printf("❌ CRON Error :: views : %v\n", err)
		return
	}

	// If there are no new views
	if len(results) == 0 {
		log.Println("✅ CRON: No new views to process.")
		return
	}

	//Update the main posts table and mark views as processed inside a transaction
	err = w.db.Transaction(func(tx *gorm.DB) error {
		// Loop through grouped results (e.g. Post A got 5 views, Post B got 12 views)
		for _, res := range results {
			// Add the new views to the existing views
			if err := tx.Model(&models.Post{}).Where("id = ?", res.PostID).
				UpdateColumn("views", gorm.Expr("views + ?", res.NewViews)).Error; err != nil {
				return err
			}
		}

		// Mark all the currently unprocessed views as processed!
		if err := tx.Model(&models.PostView{}).
			Where("processed = ?", false).
			UpdateColumn("processed", true).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("❌ CRON Error updating view counts: %v\n", err)
	} else {
		log.Printf("✅ CRON: Successfully processed views for %d posts.\n", len(results))
	}
}
