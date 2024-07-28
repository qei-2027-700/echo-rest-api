package repository

import (
	"fmt"
	"go-rest-api/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITaskRepository interface {
	// 返り値は、errorinterfaces型にしている
	// スライスのポインタ
	GetAllTasks(tasks *[]model.Task, userId uint) error
	GetTaskById(task *model.Task, userId uint, taskId uint) error
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task, userId uint, taskId uint) error
	DeleteTask(userId uint, taskId uint) error
}

// 構造体を定義
type taskRepository struct {
	db *gorm.DB
}

// コンストラクタ 受け取ったDBを使って 構造体の実態を作成する
func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &taskRepository{db}
}

// 一致するタスクの一覧を ポインタレシーバ
func (tr *taskRepository) GetAllTasks(tasks *[]model.Task, userId uint) error {
	// 引数で渡された、と一致するもの
	if err := tr.db.Joins("User").Where("user_id=?", userId).Order("created_at").Find(tasks).Error; err != nil {
		return err
	}
	return nil
}

// 指し示す先のメモリ領域に書き込む
func (tr *taskRepository) GetTaskById(task *model.Task, userId uint, taskId uint) error {
	// 引数で受け取ったuser_idと一致するものを抽出してくれる
	if err := tr.db.Joins("User").Where("user_id=?", userId).First(task, taskId).Error; err != nil {
		return err
	}
	return nil
}

// タスクのポインタを引数で渡すようにしている
func (tr *taskRepository) CreateTask(task *model.Task) error {
	if err := tr.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

func (tr *taskRepository) UpdateTask(task *model.Task, userId uint, taskId uint) error {
	// Returningのキーワードをつけると、更新しした後のObjectを、指し示す先に書き込んでくれる whereで一致するかつ、user_idが、一致するタスクに対してUPDATEする
	result := tr.db.Model(task).Clauses(clause.Returning{}).Where("id=? AND user_id=?", taskId, userId).Update("title", task.Title)

	if result.Error != nil {
		return result.Error
	}
	// 実際に更新されたレコードの数を取得する、０の場合は
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}

// 引数で渡された　一致するタスクウィおDELETEする
func (tr *taskRepository) DeleteTask(userId uint, taskId uint) error {
	result := tr.db.Where("id=? AND user_id=?", taskId, userId).Delete(&model.Task{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected < 1 {
		return fmt.Errorf("object does not exist")
	}
	return nil
}
