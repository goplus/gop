// Package driver provide Go+ "database/sql/driver" package, as "database/sql/driver" package in Go.
package driver

import (
	context "context"
	driver "database/sql/driver"
	reflect "reflect"

	gop "github.com/goplus/gop"
)

func execiColumnConverterColumnConverter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.ColumnConverter).ColumnConverter(args[1].(int))
	p.Ret(2, ret0)
}

func execiConnBegin(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(driver.Conn).Begin()
	p.Ret(1, ret0, ret1)
}

func execiConnClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Conn).Close()
	p.Ret(1, ret0)
}

func execiConnPrepare(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.Conn).Prepare(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType0(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execiConnBeginTxBeginTx(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(driver.ConnBeginTx).BeginTx(toType0(args[1]), args[2].(driver.TxOptions))
	p.Ret(3, ret0, ret1)
}

func execiConnPrepareContextPrepareContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(driver.ConnPrepareContext).PrepareContext(toType0(args[1]), args[2].(string))
	p.Ret(3, ret0, ret1)
}

func execiConnectorConnect(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.Connector).Connect(toType0(args[1]))
	p.Ret(2, ret0, ret1)
}

func execiConnectorDriver(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Connector).Driver()
	p.Ret(1, ret0)
}

func execiDriverOpen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.Driver).Open(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiDriverContextOpenConnector(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.DriverContext).OpenConnector(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiExecerExec(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(driver.Execer).Exec(args[1].(string), args[2].([]driver.Value))
	p.Ret(3, ret0, ret1)
}

func execiExecerContextExecContext(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(driver.ExecerContext).ExecContext(toType0(args[1]), args[2].(string), args[3].([]driver.NamedValue))
	p.Ret(4, ret0, ret1)
}

func execIsScanValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := driver.IsScanValue(args[0])
	p.Ret(1, ret0)
}

func execIsValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := driver.IsValue(args[0])
	p.Ret(1, ret0)
}

func execiNamedValueCheckerCheckNamedValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.NamedValueChecker).CheckNamedValue(args[1].(*driver.NamedValue))
	p.Ret(2, ret0)
}

func execmNotNullConvertValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.NotNull).ConvertValue(args[1])
	p.Ret(2, ret0, ret1)
}

func execmNullConvertValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.Null).ConvertValue(args[1])
	p.Ret(2, ret0, ret1)
}

func execiPingerPing(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.Pinger).Ping(toType0(args[1]))
	p.Ret(2, ret0)
}

func execiQueryerQuery(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(driver.Queryer).Query(args[1].(string), args[2].([]driver.Value))
	p.Ret(3, ret0, ret1)
}

func execiQueryerContextQueryContext(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(driver.QueryerContext).QueryContext(toType0(args[1]), args[2].(string), args[3].([]driver.NamedValue))
	p.Ret(4, ret0, ret1)
}

func execiResultLastInsertId(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(driver.Result).LastInsertId()
	p.Ret(1, ret0, ret1)
}

func execiResultRowsAffected(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(driver.Result).RowsAffected()
	p.Ret(1, ret0, ret1)
}

func execiRowsClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Rows).Close()
	p.Ret(1, ret0)
}

func execiRowsColumns(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Rows).Columns()
	p.Ret(1, ret0)
}

func execiRowsNext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.Rows).Next(args[1].([]driver.Value))
	p.Ret(2, ret0)
}

func execmRowsAffectedLastInsertId(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(driver.RowsAffected).LastInsertId()
	p.Ret(1, ret0, ret1)
}

func execmRowsAffectedRowsAffected(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(driver.RowsAffected).RowsAffected()
	p.Ret(1, ret0, ret1)
}

func execiRowsColumnTypeDatabaseTypeNameColumnTypeDatabaseTypeName(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.RowsColumnTypeDatabaseTypeName).ColumnTypeDatabaseTypeName(args[1].(int))
	p.Ret(2, ret0)
}

func execiRowsColumnTypeLengthColumnTypeLength(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.RowsColumnTypeLength).ColumnTypeLength(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execiRowsColumnTypeNullableColumnTypeNullable(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.RowsColumnTypeNullable).ColumnTypeNullable(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execiRowsColumnTypePrecisionScaleColumnTypePrecisionScale(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(driver.RowsColumnTypePrecisionScale).ColumnTypePrecisionScale(args[1].(int))
	p.Ret(2, ret0, ret1, ret2)
}

func execiRowsColumnTypeScanTypeColumnTypeScanType(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.RowsColumnTypeScanType).ColumnTypeScanType(args[1].(int))
	p.Ret(2, ret0)
}

func execiRowsNextResultSetHasNextResultSet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.RowsNextResultSet).HasNextResultSet()
	p.Ret(1, ret0)
}

func execiRowsNextResultSetNextResultSet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.RowsNextResultSet).NextResultSet()
	p.Ret(1, ret0)
}

func execiSessionResetterResetSession(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(driver.SessionResetter).ResetSession(toType0(args[1]))
	p.Ret(2, ret0)
}

func execiStmtClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Stmt).Close()
	p.Ret(1, ret0)
}

func execiStmtExec(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.Stmt).Exec(args[1].([]driver.Value))
	p.Ret(2, ret0, ret1)
}

func execiStmtNumInput(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Stmt).NumInput()
	p.Ret(1, ret0)
}

func execiStmtQuery(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.Stmt).Query(args[1].([]driver.Value))
	p.Ret(2, ret0, ret1)
}

func execiStmtExecContextExecContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(driver.StmtExecContext).ExecContext(toType0(args[1]), args[2].([]driver.NamedValue))
	p.Ret(3, ret0, ret1)
}

func execiStmtQueryContextQueryContext(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(driver.StmtQueryContext).QueryContext(toType0(args[1]), args[2].([]driver.NamedValue))
	p.Ret(3, ret0, ret1)
}

func execiTxCommit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Tx).Commit()
	p.Ret(1, ret0)
}

func execiTxRollback(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(driver.Tx).Rollback()
	p.Ret(1, ret0)
}

func execiValueConverterConvertValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(driver.ValueConverter).ConvertValue(args[1])
	p.Ret(2, ret0, ret1)
}

func execiValuerValue(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(driver.Valuer).Value()
	p.Ret(1, ret0, ret1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("database/sql/driver")

func init() {
	I.RegisterFuncs(
		I.Func("(ColumnConverter).ColumnConverter", (driver.ColumnConverter).ColumnConverter, execiColumnConverterColumnConverter),
		I.Func("(Conn).Begin", (driver.Conn).Begin, execiConnBegin),
		I.Func("(Conn).Close", (driver.Conn).Close, execiConnClose),
		I.Func("(Conn).Prepare", (driver.Conn).Prepare, execiConnPrepare),
		I.Func("(ConnBeginTx).BeginTx", (driver.ConnBeginTx).BeginTx, execiConnBeginTxBeginTx),
		I.Func("(ConnPrepareContext).PrepareContext", (driver.ConnPrepareContext).PrepareContext, execiConnPrepareContextPrepareContext),
		I.Func("(Connector).Connect", (driver.Connector).Connect, execiConnectorConnect),
		I.Func("(Connector).Driver", (driver.Connector).Driver, execiConnectorDriver),
		I.Func("(Driver).Open", (driver.Driver).Open, execiDriverOpen),
		I.Func("(DriverContext).OpenConnector", (driver.DriverContext).OpenConnector, execiDriverContextOpenConnector),
		I.Func("(Execer).Exec", (driver.Execer).Exec, execiExecerExec),
		I.Func("(ExecerContext).ExecContext", (driver.ExecerContext).ExecContext, execiExecerContextExecContext),
		I.Func("IsScanValue", driver.IsScanValue, execIsScanValue),
		I.Func("IsValue", driver.IsValue, execIsValue),
		I.Func("(NamedValueChecker).CheckNamedValue", (driver.NamedValueChecker).CheckNamedValue, execiNamedValueCheckerCheckNamedValue),
		I.Func("(NotNull).ConvertValue", (driver.NotNull).ConvertValue, execmNotNullConvertValue),
		I.Func("(Null).ConvertValue", (driver.Null).ConvertValue, execmNullConvertValue),
		I.Func("(Pinger).Ping", (driver.Pinger).Ping, execiPingerPing),
		I.Func("(Queryer).Query", (driver.Queryer).Query, execiQueryerQuery),
		I.Func("(QueryerContext).QueryContext", (driver.QueryerContext).QueryContext, execiQueryerContextQueryContext),
		I.Func("(Result).LastInsertId", (driver.Result).LastInsertId, execiResultLastInsertId),
		I.Func("(Result).RowsAffected", (driver.Result).RowsAffected, execiResultRowsAffected),
		I.Func("(Rows).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(Rows).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(Rows).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsAffected).LastInsertId", (driver.RowsAffected).LastInsertId, execmRowsAffectedLastInsertId),
		I.Func("(RowsAffected).RowsAffected", (driver.RowsAffected).RowsAffected, execmRowsAffectedRowsAffected),
		I.Func("(RowsColumnTypeDatabaseTypeName).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(RowsColumnTypeDatabaseTypeName).ColumnTypeDatabaseTypeName", (driver.RowsColumnTypeDatabaseTypeName).ColumnTypeDatabaseTypeName, execiRowsColumnTypeDatabaseTypeNameColumnTypeDatabaseTypeName),
		I.Func("(RowsColumnTypeDatabaseTypeName).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(RowsColumnTypeDatabaseTypeName).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsColumnTypeLength).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(RowsColumnTypeLength).ColumnTypeLength", (driver.RowsColumnTypeLength).ColumnTypeLength, execiRowsColumnTypeLengthColumnTypeLength),
		I.Func("(RowsColumnTypeLength).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(RowsColumnTypeLength).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsColumnTypeNullable).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(RowsColumnTypeNullable).ColumnTypeNullable", (driver.RowsColumnTypeNullable).ColumnTypeNullable, execiRowsColumnTypeNullableColumnTypeNullable),
		I.Func("(RowsColumnTypeNullable).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(RowsColumnTypeNullable).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsColumnTypePrecisionScale).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(RowsColumnTypePrecisionScale).ColumnTypePrecisionScale", (driver.RowsColumnTypePrecisionScale).ColumnTypePrecisionScale, execiRowsColumnTypePrecisionScaleColumnTypePrecisionScale),
		I.Func("(RowsColumnTypePrecisionScale).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(RowsColumnTypePrecisionScale).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsColumnTypeScanType).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(RowsColumnTypeScanType).ColumnTypeScanType", (driver.RowsColumnTypeScanType).ColumnTypeScanType, execiRowsColumnTypeScanTypeColumnTypeScanType),
		I.Func("(RowsColumnTypeScanType).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(RowsColumnTypeScanType).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsNextResultSet).Close", (driver.Rows).Close, execiRowsClose),
		I.Func("(RowsNextResultSet).Columns", (driver.Rows).Columns, execiRowsColumns),
		I.Func("(RowsNextResultSet).HasNextResultSet", (driver.RowsNextResultSet).HasNextResultSet, execiRowsNextResultSetHasNextResultSet),
		I.Func("(RowsNextResultSet).Next", (driver.Rows).Next, execiRowsNext),
		I.Func("(RowsNextResultSet).NextResultSet", (driver.RowsNextResultSet).NextResultSet, execiRowsNextResultSetNextResultSet),
		I.Func("(SessionResetter).ResetSession", (driver.SessionResetter).ResetSession, execiSessionResetterResetSession),
		I.Func("(Stmt).Close", (driver.Stmt).Close, execiStmtClose),
		I.Func("(Stmt).Exec", (driver.Stmt).Exec, execiStmtExec),
		I.Func("(Stmt).NumInput", (driver.Stmt).NumInput, execiStmtNumInput),
		I.Func("(Stmt).Query", (driver.Stmt).Query, execiStmtQuery),
		I.Func("(StmtExecContext).ExecContext", (driver.StmtExecContext).ExecContext, execiStmtExecContextExecContext),
		I.Func("(StmtQueryContext).QueryContext", (driver.StmtQueryContext).QueryContext, execiStmtQueryContextQueryContext),
		I.Func("(Tx).Commit", (driver.Tx).Commit, execiTxCommit),
		I.Func("(Tx).Rollback", (driver.Tx).Rollback, execiTxRollback),
		I.Func("(ValueConverter).ConvertValue", (driver.ValueConverter).ConvertValue, execiValueConverterConvertValue),
		I.Func("(Valuer).Value", (driver.Valuer).Value, execiValuerValue),
	)
	I.RegisterVars(
		I.Var("Bool", &driver.Bool),
		I.Var("DefaultParameterConverter", &driver.DefaultParameterConverter),
		I.Var("ErrBadConn", &driver.ErrBadConn),
		I.Var("ErrRemoveArgument", &driver.ErrRemoveArgument),
		I.Var("ErrSkip", &driver.ErrSkip),
		I.Var("Int32", &driver.Int32),
		I.Var("ResultNoRows", &driver.ResultNoRows),
		I.Var("String", &driver.String),
	)
	I.RegisterTypes(
		I.Type("ColumnConverter", reflect.TypeOf((*driver.ColumnConverter)(nil)).Elem()),
		I.Type("Conn", reflect.TypeOf((*driver.Conn)(nil)).Elem()),
		I.Type("ConnBeginTx", reflect.TypeOf((*driver.ConnBeginTx)(nil)).Elem()),
		I.Type("ConnPrepareContext", reflect.TypeOf((*driver.ConnPrepareContext)(nil)).Elem()),
		I.Type("Connector", reflect.TypeOf((*driver.Connector)(nil)).Elem()),
		I.Type("Driver", reflect.TypeOf((*driver.Driver)(nil)).Elem()),
		I.Type("DriverContext", reflect.TypeOf((*driver.DriverContext)(nil)).Elem()),
		I.Type("Execer", reflect.TypeOf((*driver.Execer)(nil)).Elem()),
		I.Type("ExecerContext", reflect.TypeOf((*driver.ExecerContext)(nil)).Elem()),
		I.Type("IsolationLevel", reflect.TypeOf((*driver.IsolationLevel)(nil)).Elem()),
		I.Type("NamedValue", reflect.TypeOf((*driver.NamedValue)(nil)).Elem()),
		I.Type("NamedValueChecker", reflect.TypeOf((*driver.NamedValueChecker)(nil)).Elem()),
		I.Type("NotNull", reflect.TypeOf((*driver.NotNull)(nil)).Elem()),
		I.Type("Null", reflect.TypeOf((*driver.Null)(nil)).Elem()),
		I.Type("Pinger", reflect.TypeOf((*driver.Pinger)(nil)).Elem()),
		I.Type("Queryer", reflect.TypeOf((*driver.Queryer)(nil)).Elem()),
		I.Type("QueryerContext", reflect.TypeOf((*driver.QueryerContext)(nil)).Elem()),
		I.Type("Result", reflect.TypeOf((*driver.Result)(nil)).Elem()),
		I.Type("Rows", reflect.TypeOf((*driver.Rows)(nil)).Elem()),
		I.Type("RowsAffected", reflect.TypeOf((*driver.RowsAffected)(nil)).Elem()),
		I.Type("RowsColumnTypeDatabaseTypeName", reflect.TypeOf((*driver.RowsColumnTypeDatabaseTypeName)(nil)).Elem()),
		I.Type("RowsColumnTypeLength", reflect.TypeOf((*driver.RowsColumnTypeLength)(nil)).Elem()),
		I.Type("RowsColumnTypeNullable", reflect.TypeOf((*driver.RowsColumnTypeNullable)(nil)).Elem()),
		I.Type("RowsColumnTypePrecisionScale", reflect.TypeOf((*driver.RowsColumnTypePrecisionScale)(nil)).Elem()),
		I.Type("RowsColumnTypeScanType", reflect.TypeOf((*driver.RowsColumnTypeScanType)(nil)).Elem()),
		I.Type("RowsNextResultSet", reflect.TypeOf((*driver.RowsNextResultSet)(nil)).Elem()),
		I.Type("SessionResetter", reflect.TypeOf((*driver.SessionResetter)(nil)).Elem()),
		I.Type("Stmt", reflect.TypeOf((*driver.Stmt)(nil)).Elem()),
		I.Type("StmtExecContext", reflect.TypeOf((*driver.StmtExecContext)(nil)).Elem()),
		I.Type("StmtQueryContext", reflect.TypeOf((*driver.StmtQueryContext)(nil)).Elem()),
		I.Type("Tx", reflect.TypeOf((*driver.Tx)(nil)).Elem()),
		I.Type("TxOptions", reflect.TypeOf((*driver.TxOptions)(nil)).Elem()),
		I.Type("Value", reflect.TypeOf((*driver.Value)(nil)).Elem()),
		I.Type("ValueConverter", reflect.TypeOf((*driver.ValueConverter)(nil)).Elem()),
		I.Type("Valuer", reflect.TypeOf((*driver.Valuer)(nil)).Elem()),
	)
}
