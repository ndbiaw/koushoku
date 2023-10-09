// Code generated by SQLBoiler 4.11.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Circle is an object representing the database table.
type Circle struct {
	ID   int64  `boil:"id" json:"id" toml:"id" yaml:"id"`
	Slug string `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`

	R *circleR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L circleL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CircleColumns = struct {
	ID   string
	Slug string
	Name string
}{
	ID:   "id",
	Slug: "slug",
	Name: "name",
}

var CircleTableColumns = struct {
	ID   string
	Slug string
	Name string
}{
	ID:   "circle.id",
	Slug: "circle.slug",
	Name: "circle.name",
}

// Generated where

var CircleWhere = struct {
	ID   whereHelperint64
	Slug whereHelperstring
	Name whereHelperstring
}{
	ID:   whereHelperint64{field: "\"circle\".\"id\""},
	Slug: whereHelperstring{field: "\"circle\".\"slug\""},
	Name: whereHelperstring{field: "\"circle\".\"name\""},
}

// CircleRels is where relationship names are stored.
var CircleRels = struct {
	Archives string
}{
	Archives: "Archives",
}

// circleR is where relationships are stored.
type circleR struct {
	Archives ArchiveSlice `boil:"Archives" json:"Archives" toml:"Archives" yaml:"Archives"`
}

// NewStruct creates a new relationship struct
func (*circleR) NewStruct() *circleR {
	return &circleR{}
}

func (r *circleR) GetArchives() ArchiveSlice {
	if r == nil {
		return nil
	}
	return r.Archives
}

// circleL is where Load methods for each relationship are stored.
type circleL struct{}

var (
	circleAllColumns            = []string{"id", "slug", "name"}
	circleColumnsWithoutDefault = []string{}
	circleColumnsWithDefault    = []string{"id", "slug", "name"}
	circlePrimaryKeyColumns     = []string{"id"}
	circleGeneratedColumns      = []string{}
)

type (
	// CircleSlice is an alias for a slice of pointers to Circle.
	// This should almost always be used instead of []Circle.
	CircleSlice []*Circle
	// CircleHook is the signature for custom Circle hook methods
	CircleHook func(boil.Executor, *Circle) error

	circleQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	circleType                 = reflect.TypeOf(&Circle{})
	circleMapping              = queries.MakeStructMapping(circleType)
	circlePrimaryKeyMapping, _ = queries.BindMapping(circleType, circleMapping, circlePrimaryKeyColumns)
	circleInsertCacheMut       sync.RWMutex
	circleInsertCache          = make(map[string]insertCache)
	circleUpdateCacheMut       sync.RWMutex
	circleUpdateCache          = make(map[string]updateCache)
	circleUpsertCacheMut       sync.RWMutex
	circleUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var circleAfterSelectHooks []CircleHook

var circleBeforeInsertHooks []CircleHook
var circleAfterInsertHooks []CircleHook

var circleBeforeUpdateHooks []CircleHook
var circleAfterUpdateHooks []CircleHook

var circleBeforeDeleteHooks []CircleHook
var circleAfterDeleteHooks []CircleHook

var circleBeforeUpsertHooks []CircleHook
var circleAfterUpsertHooks []CircleHook

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Circle) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range circleAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Circle) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range circleBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Circle) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range circleAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Circle) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range circleBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Circle) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range circleAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Circle) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range circleBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Circle) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range circleAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Circle) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range circleBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Circle) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range circleAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddCircleHook registers your hook function for all future operations.
func AddCircleHook(hookPoint boil.HookPoint, circleHook CircleHook) {
	switch hookPoint {
	case boil.AfterSelectHook:
		circleAfterSelectHooks = append(circleAfterSelectHooks, circleHook)
	case boil.BeforeInsertHook:
		circleBeforeInsertHooks = append(circleBeforeInsertHooks, circleHook)
	case boil.AfterInsertHook:
		circleAfterInsertHooks = append(circleAfterInsertHooks, circleHook)
	case boil.BeforeUpdateHook:
		circleBeforeUpdateHooks = append(circleBeforeUpdateHooks, circleHook)
	case boil.AfterUpdateHook:
		circleAfterUpdateHooks = append(circleAfterUpdateHooks, circleHook)
	case boil.BeforeDeleteHook:
		circleBeforeDeleteHooks = append(circleBeforeDeleteHooks, circleHook)
	case boil.AfterDeleteHook:
		circleAfterDeleteHooks = append(circleAfterDeleteHooks, circleHook)
	case boil.BeforeUpsertHook:
		circleBeforeUpsertHooks = append(circleBeforeUpsertHooks, circleHook)
	case boil.AfterUpsertHook:
		circleAfterUpsertHooks = append(circleAfterUpsertHooks, circleHook)
	}
}

// OneG returns a single circle record from the query using the global executor.
func (q circleQuery) OneG() (*Circle, error) {
	return q.One(boil.GetDB())
}

// One returns a single circle record from the query.
func (q circleQuery) One(exec boil.Executor) (*Circle, error) {
	o := &Circle{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(nil, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for circle")
	}

	if err := o.doAfterSelectHooks(exec); err != nil {
		return o, err
	}

	return o, nil
}

// AllG returns all Circle records from the query using the global executor.
func (q circleQuery) AllG() (CircleSlice, error) {
	return q.All(boil.GetDB())
}

// All returns all Circle records from the query.
func (q circleQuery) All(exec boil.Executor) (CircleSlice, error) {
	var o []*Circle

	err := q.Bind(nil, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Circle slice")
	}

	if len(circleAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountG returns the count of all Circle records in the query using the global executor
func (q circleQuery) CountG() (int64, error) {
	return q.Count(boil.GetDB())
}

// Count returns the count of all Circle records in the query.
func (q circleQuery) Count(exec boil.Executor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count circle rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q circleQuery) ExistsG() (bool, error) {
	return q.Exists(boil.GetDB())
}

// Exists checks if the row exists in the table.
func (q circleQuery) Exists(exec boil.Executor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow(exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if circle exists")
	}

	return count > 0, nil
}

// Archives retrieves all the archive's Archives with an executor.
func (o *Circle) Archives(mods ...qm.QueryMod) archiveQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.InnerJoin("\"archive_circles\" on \"archive\".\"id\" = \"archive_circles\".\"archive_id\""),
		qm.Where("\"archive_circles\".\"circle_id\"=?", o.ID),
	)

	return Archives(queryMods...)
}

// LoadArchives allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (circleL) LoadArchives(e boil.Executor, singular bool, maybeCircle interface{}, mods queries.Applicator) error {
	var slice []*Circle
	var object *Circle

	if singular {
		object = maybeCircle.(*Circle)
	} else {
		slice = *maybeCircle.(*[]*Circle)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &circleR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &circleR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.Select("\"archive\".\"id\", \"archive\".\"path\", \"archive\".\"created_at\", \"archive\".\"updated_at\", \"archive\".\"published_at\", \"archive\".\"title\", \"archive\".\"slug\", \"archive\".\"pages\", \"archive\".\"size\", \"archive\".\"expunged\", \"archive\".\"source\", \"archive\".\"submission_id\", \"archive\".\"redirect_id\", \"a\".\"circle_id\""),
		qm.From("\"archive\""),
		qm.InnerJoin("\"archive_circles\" as \"a\" on \"archive\".\"id\" = \"a\".\"archive_id\""),
		qm.WhereIn("\"a\".\"circle_id\" in ?", args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.Query(e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load archive")
	}

	var resultSlice []*Archive

	var localJoinCols []int64
	for results.Next() {
		one := new(Archive)
		var localJoinCol int64

		err = results.Scan(&one.ID, &one.Path, &one.CreatedAt, &one.UpdatedAt, &one.PublishedAt, &one.Title, &one.Slug, &one.Pages, &one.Size, &one.Expunged, &one.Source, &one.SubmissionID, &one.RedirectID, &localJoinCol)
		if err != nil {
			return errors.Wrap(err, "failed to scan eager loaded results for archive")
		}
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice archive")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on archive")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for archive")
	}

	if len(archiveAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Archives = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &archiveR{}
			}
			foreign.R.Circles = append(foreign.R.Circles, object)
		}
		return nil
	}

	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			if local.ID == localJoinCol {
				local.R.Archives = append(local.R.Archives, foreign)
				if foreign.R == nil {
					foreign.R = &archiveR{}
				}
				foreign.R.Circles = append(foreign.R.Circles, local)
				break
			}
		}
	}

	return nil
}

// AddArchivesG adds the given related objects to the existing relationships
// of the circle, optionally inserting them as new records.
// Appends related to o.R.Archives.
// Sets related.R.Circles appropriately.
// Uses the global database handle.
func (o *Circle) AddArchivesG(insert bool, related ...*Archive) error {
	return o.AddArchives(boil.GetDB(), insert, related...)
}

// AddArchives adds the given related objects to the existing relationships
// of the circle, optionally inserting them as new records.
// Appends related to o.R.Archives.
// Sets related.R.Circles appropriately.
func (o *Circle) AddArchives(exec boil.Executor, insert bool, related ...*Archive) error {
	var err error
	for _, rel := range related {
		if insert {
			if err = rel.Insert(exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}
	}

	for _, rel := range related {
		query := "insert into \"archive_circles\" (\"circle_id\", \"archive_id\") values ($1, $2)"
		values := []interface{}{o.ID, rel.ID}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}
		_, err = exec.Exec(query, values...)
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	if o.R == nil {
		o.R = &circleR{
			Archives: related,
		}
	} else {
		o.R.Archives = append(o.R.Archives, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &archiveR{
				Circles: CircleSlice{o},
			}
		} else {
			rel.R.Circles = append(rel.R.Circles, o)
		}
	}
	return nil
}

// SetArchivesG removes all previously related items of the
// circle replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Circles's Archives accordingly.
// Replaces o.R.Archives with related.
// Sets related.R.Circles's Archives accordingly.
// Uses the global database handle.
func (o *Circle) SetArchivesG(insert bool, related ...*Archive) error {
	return o.SetArchives(boil.GetDB(), insert, related...)
}

// SetArchives removes all previously related items of the
// circle replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Circles's Archives accordingly.
// Replaces o.R.Archives with related.
// Sets related.R.Circles's Archives accordingly.
func (o *Circle) SetArchives(exec boil.Executor, insert bool, related ...*Archive) error {
	query := "delete from \"archive_circles\" where \"circle_id\" = $1"
	values := []interface{}{o.ID}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	removeArchivesFromCirclesSlice(o, related)
	if o.R != nil {
		o.R.Archives = nil
	}

	return o.AddArchives(exec, insert, related...)
}

// RemoveArchivesG relationships from objects passed in.
// Removes related items from R.Archives (uses pointer comparison, removal does not keep order)
// Sets related.R.Circles.
// Uses the global database handle.
func (o *Circle) RemoveArchivesG(related ...*Archive) error {
	return o.RemoveArchives(boil.GetDB(), related...)
}

// RemoveArchives relationships from objects passed in.
// Removes related items from R.Archives (uses pointer comparison, removal does not keep order)
// Sets related.R.Circles.
func (o *Circle) RemoveArchives(exec boil.Executor, related ...*Archive) error {
	if len(related) == 0 {
		return nil
	}

	var err error
	query := fmt.Sprintf(
		"delete from \"archive_circles\" where \"circle_id\" = $1 and \"archive_id\" in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, len(related), 2, 1),
	)
	values := []interface{}{o.ID}
	for _, rel := range related {
		values = append(values, rel.ID)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err = exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	removeArchivesFromCirclesSlice(o, related)
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Archives {
			if rel != ri {
				continue
			}

			ln := len(o.R.Archives)
			if ln > 1 && i < ln-1 {
				o.R.Archives[i] = o.R.Archives[ln-1]
			}
			o.R.Archives = o.R.Archives[:ln-1]
			break
		}
	}

	return nil
}

func removeArchivesFromCirclesSlice(o *Circle, related []*Archive) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.Circles {
			if o.ID != ri.ID {
				continue
			}

			ln := len(rel.R.Circles)
			if ln > 1 && i < ln-1 {
				rel.R.Circles[i] = rel.R.Circles[ln-1]
			}
			rel.R.Circles = rel.R.Circles[:ln-1]
			break
		}
	}
}

// Circles retrieves all the records using an executor.
func Circles(mods ...qm.QueryMod) circleQuery {
	mods = append(mods, qm.From("\"circle\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"circle\".*"})
	}

	return circleQuery{q}
}

// FindCircleG retrieves a single record by ID.
func FindCircleG(iD int64, selectCols ...string) (*Circle, error) {
	return FindCircle(boil.GetDB(), iD, selectCols...)
}

// FindCircle retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCircle(exec boil.Executor, iD int64, selectCols ...string) (*Circle, error) {
	circleObj := &Circle{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"circle\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(nil, exec, circleObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from circle")
	}

	if err = circleObj.doAfterSelectHooks(exec); err != nil {
		return circleObj, err
	}

	return circleObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Circle) InsertG(columns boil.Columns) error {
	return o.Insert(boil.GetDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Circle) Insert(exec boil.Executor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no circle provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(circleColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	circleInsertCacheMut.RLock()
	cache, cached := circleInsertCache[key]
	circleInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			circleAllColumns,
			circleColumnsWithDefault,
			circleColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(circleType, circleMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(circleType, circleMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"circle\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"circle\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into circle")
	}

	if !cached {
		circleInsertCacheMut.Lock()
		circleInsertCache[key] = cache
		circleInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Circle record using the global executor.
// See Update for more documentation.
func (o *Circle) UpdateG(columns boil.Columns) error {
	return o.Update(boil.GetDB(), columns)
}

// Update uses an executor to update the Circle.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Circle) Update(exec boil.Executor, columns boil.Columns) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(columns, nil)
	circleUpdateCacheMut.RLock()
	cache, cached := circleUpdateCache[key]
	circleUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			circleAllColumns,
			circlePrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update circle, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"circle\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, circlePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(circleType, circleMapping, append(wl, circlePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update circle row")
	}

	if !cached {
		circleUpdateCacheMut.Lock()
		circleUpdateCache[key] = cache
		circleUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllG updates all rows with the specified column values.
func (q circleQuery) UpdateAllG(cols M) error {
	return q.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q circleQuery) UpdateAll(exec boil.Executor, cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for circle")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o CircleSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CircleSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), circlePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"circle\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, circlePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in circle slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Circle) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Circle) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no circle provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(circleColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	circleUpsertCacheMut.RLock()
	cache, cached := circleUpsertCache[key]
	circleUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			circleAllColumns,
			circleColumnsWithDefault,
			circleColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			circleAllColumns,
			circlePrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert circle, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(circlePrimaryKeyColumns))
			copy(conflict, circlePrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"circle\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(circleType, circleMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(circleType, circleMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert circle")
	}

	if !cached {
		circleUpsertCacheMut.Lock()
		circleUpsertCache[key] = cache
		circleUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteG deletes a single Circle record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Circle) DeleteG() error {
	return o.Delete(boil.GetDB())
}

// Delete deletes a single Circle record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Circle) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Circle provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), circlePrimaryKeyMapping)
	sql := "DELETE FROM \"circle\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from circle")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

func (q circleQuery) DeleteAllG() error {
	return q.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all matching rows.
func (q circleQuery) DeleteAll(exec boil.Executor) error {
	if q.Query == nil {
		return errors.New("models: no circleQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec(exec)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from circle")
	}

	return nil
}

// DeleteAllG deletes all rows in the slice.
func (o CircleSlice) DeleteAllG() error {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CircleSlice) DeleteAll(exec boil.Executor) error {
	if len(o) == 0 {
		return nil
	}

	if len(circleBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), circlePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"circle\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, circlePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from circle slice")
	}

	if len(circleAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Circle) ReloadG() error {
	if o == nil {
		return errors.New("models: no Circle provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Circle) Reload(exec boil.Executor) error {
	ret, err := FindCircle(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CircleSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty CircleSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CircleSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CircleSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), circlePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"circle\".* FROM \"circle\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, circlePrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(nil, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in CircleSlice")
	}

	*o = slice

	return nil
}

// CircleExistsG checks if the Circle row exists.
func CircleExistsG(iD int64) (bool, error) {
	return CircleExists(boil.GetDB(), iD)
}

// CircleExists checks if the Circle row exists.
func CircleExists(exec boil.Executor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"circle\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, iD)
	}
	row := exec.QueryRow(sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if circle exists")
	}

	return exists, nil
}