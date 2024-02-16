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

func (qb *viewDateManager) GetManyViewDates(ctx context.Context, ids []int) ([][]time.Time, error) {
	return qb.tableMgr.getManyDates(ctx, ids)
}

func (qb *viewDateManager) CountViews(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.getCount(ctx, id)
}

func (qb *viewDateManager) GetManyViewCount(ctx context.Context, ids []int) ([]int, error) {
	return qb.tableMgr.getManyCount(ctx, ids)
}

func (qb *viewDateManager) CountAllViews(ctx context.Context) (int, error) {
	return qb.tableMgr.getAllCount(ctx)
}

func (qb *viewDateManager) CountUniqueViews(ctx context.Context) (int, error) {
	return qb.tableMgr.getUniqueCount(ctx)
}

func (qb *viewDateManager) LastView(ctx context.Context, id int) (*time.Time, error) {
	return qb.tableMgr.getLastDate(ctx, id)
}

func (qb *viewDateManager) GetManyLastViewed(ctx context.Context, ids []int) ([]*time.Time, error) {
	return qb.tableMgr.getManyLastDate(ctx, ids)

}

func (qb *viewDateManager) AddViews(ctx context.Context, id int, dates []time.Time) ([]time.Time, error) {
	return qb.tableMgr.addDates(ctx, id, dates)
}

func (qb *viewDateManager) DeleteViews(ctx context.Context, id int, dates []time.Time) ([]time.Time, error) {
	return qb.tableMgr.deleteDates(ctx, id, dates)
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

func (qb *oDateManager) GetManyODates(ctx context.Context, ids []int) ([][]time.Time, error) {
	return qb.tableMgr.getManyDates(ctx, ids)
}

func (qb *oDateManager) GetOCount(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.getCount(ctx, id)
}

func (qb *oDateManager) GetManyOCount(ctx context.Context, ids []int) ([]int, error) {
	return qb.tableMgr.getManyCount(ctx, ids)
}

func (qb *oDateManager) GetAllOCount(ctx context.Context) (int, error) {
	return qb.tableMgr.getAllCount(ctx)
}

func (qb *oDateManager) GetUniqueOCount(ctx context.Context) (int, error) {
	return qb.tableMgr.getUniqueCount(ctx)
}

func (qb *oDateManager) AddO(ctx context.Context, id int, dates []time.Time) ([]time.Time, error) {
	return qb.tableMgr.addDates(ctx, id, dates)
}

func (qb *oDateManager) DeleteO(ctx context.Context, id int, dates []time.Time) ([]time.Time, error) {
	return qb.tableMgr.deleteDates(ctx, id, dates)
}

func (qb *oDateManager) ResetO(ctx context.Context, id int) (int, error) {
	return qb.tableMgr.deleteAllDates(ctx, id)
}
