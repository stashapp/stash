package sqlite

import (
	"context"
	"time"
)

type viewDateManager struct {
	tableMgr *viewHistoryTable
}

func (qb *viewDateManager) GetViewDates(ctx context.Context, id int) ([]time.Time, error) {
	return qb.tableMgr.getDates(ctx, id)
}

func (qb *viewDateManager) CountViews(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.getCount(ctx, id)
}

func (qb *viewDateManager) LastView(ctx context.Context, id int) (*time.Time, error) {
	return qb.tableMgr.getLastDate(ctx, id)
}

func (qb *viewDateManager) AddView(ctx context.Context, id int, date *time.Time) (int, error) {
	return qb.tableMgr.addDate(ctx, id, date)
}

func (qb *viewDateManager) DeleteView(ctx context.Context, id int, date *time.Time) (int, error) {
	return qb.tableMgr.deleteDate(ctx, id, date)
}

func (qb *viewDateManager) DeleteAllViews(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.deleteAllDates(ctx, id)
}

type oDateManager struct {
	tableMgr *viewHistoryTable
}

func (qb *oDateManager) GetODates(ctx context.Context, id int) ([]time.Time, error) {
	return qb.tableMgr.getDates(ctx, id)
}

func (qb *oDateManager) GetOCount(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.getCount(ctx, id)
}

func (qb *oDateManager) AddO(ctx context.Context, id int, date *time.Time) (int, error) {
	return qb.tableMgr.addDate(ctx, id, date)
}

func (qb *oDateManager) DeleteO(ctx context.Context, id int, date *time.Time) (int, error) {
	return qb.tableMgr.deleteDate(ctx, id, date)
}

func (qb *oDateManager) ResetO(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.deleteAllDates(ctx, id)
}
