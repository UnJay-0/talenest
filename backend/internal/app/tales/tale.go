package tales

import (
	"errors"
	"fmt"
	"talenest/backend/internal/app/status"
	"talenest/backend/internal/app/tags"
	"talenest/backend/internal/utils"
	"time"
)

type Tale struct {
	Id       int
	Name     string
	Summary  string
	ParentId int
	Status   status.Status
	Tags     []tags.Tag
	created  time.Time
	updated  time.Time
	deleted  time.Time
}

func Create() (tale *Tale) {
	return &Tale{
		ParentId: 1, // Root Parent
		Status:   status.GetDefault(),
		created:  time.Now(),
		updated:  time.Now(),
	}
}

func (tale *Tale) Update() {
	tale.updated = time.Now()
}

func (tale *Tale) Delete() {
	tale.deleted = time.Now()
}

func (tale *Tale) setCreated(datetime string) error {
	created, err := time.Parse(utils.DATETIME_FORMAT, datetime)
	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to parse time value")
	}
	tale.created = created
	return nil
}

func (tale *Tale) setUpdated(datetime string) error {
	updated, err := time.Parse(utils.DATETIME_FORMAT, datetime)
	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to parse time value")
	}
	tale.updated = updated
	return nil
}

func (tale *Tale) setDeleted(datetime string) error {
	if datetime == "" {
		return nil
	}
	deleted, err := time.Parse(utils.DATETIME_FORMAT, datetime)
	if err != nil {
		return errors.New("Failed to parse time value")
	}
	tale.deleted = deleted
	return nil
}

func (tale *Tale) String() string {
	return fmt.Sprintf("Tale [%d]: %s (%s)\nparent: %d\nstatus:%s\ncreated: %s\nupdated: %s\n",
		tale.Id,
		tale.Name,
		tale.Summary,
		tale.ParentId,
		tale.Status.Name,
		tale.created.Format(utils.DATETIME_FORMAT),
		tale.updated.Format(utils.DATETIME_FORMAT))
}
