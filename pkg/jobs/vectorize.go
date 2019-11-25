package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mlhamel/survilleray/pkg/models"
	"github.com/mlhamel/survilleray/pkg/runtime"
)

type VectorizeJob struct {
	context *runtime.Context
}

func NewVectorizeJob(context *runtime.Context) *VectorizeJob {
	return &VectorizeJob{context}
}

func (job *VectorizeJob) Run() error {
	pointRepos := models.NewPointRepository(job.context)
	vectorRepos := models.NewVectorRepository(job.context)
	points, err := pointRepos.FindByVectorizedAt(nil)

	if err != nil {
		return fmt.Errorf("Cannot find points to vectorize: %w", err)
	}

	if len(points) == 0 {
		log.Printf("Cannot find any point to vectorize")
	}

	for i := 0; i < len(points); i++ {
		tx := job.context.Database().Begin()
		if tx.Error != nil {
			return err
		}

		point := points[i]

		vector, err := vectorRepos.FindByCallSign(point.CallSign)
		if err != nil && !gorm.IsRecordNotFoundError(err) {
			return fmt.Errorf("Cannot find vector: %w", err)
		}

		if gorm.IsRecordNotFoundError(err) {
			log.Printf("Creating vector for point %s", point.String())
			vector = models.NewVectorFromPoint(&point)

			if err = job.context.Database().Create(&vector).Error; err != nil {
				return fmt.Errorf("Cannot create vector: %w", err)
			}
		}

		if err = vectorRepos.AppendPoints(vector, []models.Point{point}); err != nil {
			return fmt.Errorf("Cannot add point to the matching vector: %w", err)
		}

		if err = pointRepos.Update(&point, map[string]interface{}{"VectorizedAt": time.Now()}); err != nil {
			return fmt.Errorf("Cannot update VectorizedAt for a point: %w", err)
		}

		if point.OnGround {
			if err = vectorRepos.Update(vector, map[string]interface{}{"Closed": true}); err != nil {
				return fmt.Errorf("Cannot update Closed for a vector: %w", err)
			}
		}

		if err = tx.Commit().Error; err != nil {
			return fmt.Errorf("Cannot commit transaction: %w", err)
		}
	}

	return nil
}
