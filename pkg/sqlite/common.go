package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type oCounterManager struct {
	tableMgr *table
}
type playCounterManager struct {
	tableMgr *table
}

type oDateManager struct {
	tableMgr *table
}

type playDateManager struct {
	tableMgr *table
}

func (qb *oCounterManager) getOCounter(ctx context.Context, id int) (int, error) {
	q := dialect.From(qb.tableMgr.table).Select("o_counter").Where(goqu.Ex{"id": id})

	const single = true
	var ret int
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		if err := rows.Scan(&ret); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *playCounterManager) getPlayCount(ctx context.Context, id int) (int, error) {
	q := dialect.From(qb.tableMgr.table).Select("play_count").Where(goqu.Ex{"id": id})

	const single = true
	var ret int
	if err := queryFunc(ctx, q, single, func(rows *sqlx.Rows) error {
		if err := rows.Scan(&ret); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (qb *oDateManager) AddODate(ctx context.Context, id int) error {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return err
	}

	if err := qb.tableMgr.addODateByID(ctx, id, goqu.Record{
		"scene_id": id,
		"odate":    time.Now().Local().Format(time.RFC3339Nano),
	}); err != nil {
		return err
	}

	return nil
}

func (qb *oDateManager) DeleteODate(ctx context.Context, id int) error {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return err
	}

	if err := qb.tableMgr.deleteODateByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func (qb *oDateManager) ResetODate(ctx context.Context, sceneID int) error {
	if err := qb.tableMgr.checkIDExists(ctx, sceneID); err != nil {
		return err
	}

	if err := qb.tableMgr.resetODateByID(ctx, sceneID); err != nil {
		return err
	}

	return nil
}

func (qb *playDateManager) AddPlayDate(ctx context.Context, id int) error {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return err
	}

	if err := qb.tableMgr.addPlayDateByID(ctx, id, goqu.Record{
		"scene_id": id,
		"playdate": time.Now().Local().Format(time.RFC3339Nano),
	}); err != nil {
		return err
	}

	return nil
}

func (qb *playDateManager) DeletePlayDate(ctx context.Context, id int) error {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return err
	}

	if err := qb.tableMgr.deletePlayDateByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func (qb *playDateManager) ResetPlayDate(ctx context.Context, sceneID int) error {
	if err := qb.tableMgr.checkIDExists(ctx, sceneID); err != nil {
		return err
	}

	if err := qb.tableMgr.resetPlayDateByID(ctx, sceneID); err != nil {
		return err
	}

	return nil
}

func (qb *oCounterManager) IncrementOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": goqu.L("o_counter + 1"),
	}); err != nil {
		return 0, err
	}

	oDateMgr := &oDateManager{tableMgr: qb.tableMgr}
	if err := oDateMgr.AddODate(ctx, id); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *oCounterManager) DecrementOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	table := qb.tableMgr.table
	q := dialect.Update(table).Set(goqu.Record{
		"o_counter": goqu.L("o_counter - 1"),
	}).Where(qb.tableMgr.byID(id), goqu.L("o_counter > 0"))

	if _, err := exec(ctx, q); err != nil {
		return 0, fmt.Errorf("updating %s: %w", table.GetTable(), err)
	}

	oDateMgr := &oDateManager{tableMgr: qb.tableMgr}
	if err := oDateMgr.DeleteODate(ctx, id); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *oCounterManager) ResetOCounter(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"o_counter": 0,
	}); err != nil {
		return 0, err
	}

	oDateMgr := &oDateManager{tableMgr: qb.tableMgr}
	if err := oDateMgr.ResetODate(ctx, id); err != nil {
		return 0, err
	}

	return qb.getOCounter(ctx, id)
}

func (qb *playCounterManager) IncrementWatchCount(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"play_count":     goqu.L("play_count + 1"),
		"last_played_at": time.Now(),
	}); err != nil {
		return 0, err
	}

	playDateMgr := &playDateManager{tableMgr: qb.tableMgr}
	if err := playDateMgr.AddPlayDate(ctx, id); err != nil {
		return 0, err
	}

	return qb.getPlayCount(ctx, id)
}

func (qb *playCounterManager) DecrementWatchCount(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"play_count":     goqu.L("play_count - 1"),
		"resume_time":    0.0,
		"last_played_at": goqu.L("last_played_at == null"),
	}); err != nil {
		return 0, err
	}

	playDateMgr := &playDateManager{tableMgr: qb.tableMgr}
	if err := playDateMgr.DeletePlayDate(ctx, id); err != nil {
		return 0, err
	}

	return qb.getPlayCount(ctx, id)
}

func (qb *playCounterManager) ResetWatchCount(ctx context.Context, id int) (int, error) {
	if err := qb.tableMgr.checkIDExists(ctx, id); err != nil {
		return 0, err
	}

	if err := qb.tableMgr.updateByID(ctx, id, goqu.Record{
		"play_count":     0,
		"resume_time":    0.0,
		"last_played_at": goqu.L("last_played_at == null"),
		"play_duration":  0.0,
	}); err != nil {
		return 0, err
	}

	playDateMgr := &playDateManager{tableMgr: qb.tableMgr}
	if err := playDateMgr.ResetPlayDate(ctx, id); err != nil {
		return 0, err
	}

	return qb.getPlayCount(ctx, id)
}
