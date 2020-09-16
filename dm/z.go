/*
 * Copyright (c) 2000-2018, 达梦数据库有限公司.
 * All rights reserved.
 */

package dm

import (
	"context"
	"database/sql/driver"
	"reflect"
)

type baseFilter struct {
}

//DmDriver
func (bf *baseFilter) DmDriverOpen(filterChain *filterChain, d *DmDriver, dsn string) (*DmConnection, error) {
	return d.open(dsn)
}

func (bf *baseFilter) DmDriverOpenConnector(filterChain *filterChain, d *DmDriver, dsn string) (*DmConnector, error) {
	return d.openConnector(dsn)
}

//DmConnector
func (bf *baseFilter) DmConnectorConnect(filterChain *filterChain, c *DmConnector, ctx context.Context) (*DmConnection, error) {
	return c.connect(ctx)
}

func (bf *baseFilter) DmConnectorDriver(filterChain *filterChain, c *DmConnector) *DmDriver {
	return c.driver()
}

//DmConnection
func (bf *baseFilter) DmConnectionBegin(filterChain *filterChain, c *DmConnection) (*DmConnection, error) {
	return c.begin()
}

func (bf *baseFilter) DmConnectionBeginTx(filterChain *filterChain, c *DmConnection, ctx context.Context, opts driver.TxOptions) (*DmConnection, error) {
	return c.beginTx(ctx, opts)
}

func (bf *baseFilter) DmConnectionCommit(filterChain *filterChain, c *DmConnection) error {
	return c.commit()
}

func (bf *baseFilter) DmConnectionRollback(filterChain *filterChain, c *DmConnection) error {
	return c.rollback()
}

func (bf *baseFilter) DmConnectionClose(filterChain *filterChain, c *DmConnection) error {
	return c.close()
}

func (bf *baseFilter) DmConnectionPing(filterChain *filterChain, c *DmConnection, ctx context.Context) error {
	return c.ping(ctx)
}

func (bf *baseFilter) DmConnectionExec(filterChain *filterChain, c *DmConnection, query string, args []driver.Value) (*DmResult, error) {
	return c.exec(query, args)
}

func (bf *baseFilter) DmConnectionExecContext(filterChain *filterChain, c *DmConnection, ctx context.Context, query string, args []driver.NamedValue) (*DmResult, error) {
	return c.execContext(ctx, query, args)
}

func (bf *baseFilter) DmConnectionQuery(filterChain *filterChain, c *DmConnection, query string, args []driver.Value) (*DmRows, error) {
	return c.query(query, args)
}

func (bf *baseFilter) DmConnectionQueryContext(filterChain *filterChain, c *DmConnection, ctx context.Context, query string, args []driver.NamedValue) (*DmRows, error) {
	return c.queryContext(ctx, query, args)
}

func (bf *baseFilter) DmConnectionPrepare(filterChain *filterChain, c *DmConnection, query string) (*DmStatement, error) {
	return c.prepare(query)
}

func (bf *baseFilter) DmConnectionPrepareContext(filterChain *filterChain, c *DmConnection, ctx context.Context, query string) (*DmStatement, error) {
	return c.prepareContext(ctx, query)
}

func (bf *baseFilter) DmConnectionResetSession(filterChain *filterChain, c *DmConnection, ctx context.Context) error {
	return c.resetSession(ctx)
}

func (bf *baseFilter) DmConnectionCheckNamedValue(filterChain *filterChain, c *DmConnection, nv *driver.NamedValue) error {
	return c.checkNamedValue(nv)
}

//DmStatement
func (bf *baseFilter) DmStatementClose(filterChain *filterChain, s *DmStatement) error {
	return s.close()
}

func (bf *baseFilter) DmStatementNumInput(filterChain *filterChain, s *DmStatement) int {
	return s.numInput()
}

func (bf *baseFilter) DmStatementExec(filterChain *filterChain, s *DmStatement, args []driver.Value) (*DmResult, error) {
	return s.exec(args)
}

func (bf *baseFilter) DmStatementExecContext(filterChain *filterChain, s *DmStatement, ctx context.Context, args []driver.NamedValue) (*DmResult, error) {
	return s.execContext(ctx, args)
}

func (bf *baseFilter) DmStatementQuery(filterChain *filterChain, s *DmStatement, args []driver.Value) (*DmRows, error) {
	return s.query(args)
}

func (bf *baseFilter) DmStatementQueryContext(filterChain *filterChain, s *DmStatement, ctx context.Context, args []driver.NamedValue) (*DmRows, error) {
	return s.queryContext(ctx, args)
}

func (bf *baseFilter) DmStatementCheckNamedValue(filterChain *filterChain, s *DmStatement, nv *driver.NamedValue) error {
	return s.checkNamedValue(nv)
}

//DmResult
func (bf *baseFilter) DmResultLastInsertId(filterChain *filterChain, r *DmResult) (int64, error) {
	return r.lastInsertId()
}

func (bf *baseFilter) DmResultRowsAffected(filterChain *filterChain, r *DmResult) (int64, error) {
	return r.rowsAffected()
}

//DmRows
func (bf *baseFilter) DmRowsColumns(filterChain *filterChain, r *DmRows) []string {
	return r.columns()
}

func (bf *baseFilter) DmRowsClose(filterChain *filterChain, r *DmRows) error {
	return r.close()
}

func (bf *baseFilter) DmRowsNext(filterChain *filterChain, r *DmRows, dest []driver.Value) error {
	return r.next(dest)
}

func (bf *baseFilter) DmRowsHasNextResultSet(filterChain *filterChain, r *DmRows) bool {
	return r.hasNextResultSet()
}

func (bf *baseFilter) DmRowsNextResultSet(filterChain *filterChain, r *DmRows) error {
	return r.nextResultSet()
}

func (bf *baseFilter) DmRowsColumnTypeScanType(filterChain *filterChain, r *DmRows, index int) reflect.Type {
	return r.columnTypeScanType(index)
}

func (bf *baseFilter) DmRowsColumnTypeDatabaseTypeName(filterChain *filterChain, r *DmRows, index int) string {
	return r.columnTypeDatabaseTypeName(index)
}

func (bf *baseFilter) DmRowsColumnTypeLength(filterChain *filterChain, r *DmRows, index int) (length int64, ok bool) {
	return r.columnTypeLength(index)
}

func (bf *baseFilter) DmRowsColumnTypeNullable(filterChain *filterChain, r *DmRows, index int) (nullable, ok bool) {
	return r.columnTypeNullable(index)
}

func (bf *baseFilter) DmRowsColumnTypePrecisionScale(filterChain *filterChain, r *DmRows, index int) (precision, scale int64, ok bool) {
	return r.columnTypePrecisionScale(index)
}
