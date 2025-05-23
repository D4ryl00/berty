package sqlite

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Migrator struct {
	migrator.Migrator
}

func (m *Migrator) RunWithoutForeignKey(fc func() error) error {
	var enabled int
	m.DB.Raw("PRAGMA foreign_keys").Scan(&enabled)
	if enabled == 1 {
		m.DB.Exec("PRAGMA foreign_keys = OFF")
		defer m.DB.Exec("PRAGMA foreign_keys = ON")
	}

	return fc()
}

func (m Migrator) HasTable(value interface{}) bool {
	var count int
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", stmt.Table).Row().Scan(&count)
	})
	return count > 0
}

func (m Migrator) DropTable(values ...interface{}) error {
	return m.RunWithoutForeignKey(func() error {
		values = m.ReorderModels(values, false)
		tx := m.DB.Session(&gorm.Session{})

		for i := len(values) - 1; i >= 0; i-- {
			if err := m.RunWithValue(values[i], func(stmt *gorm.Statement) error {
				return tx.Exec("DROP TABLE IF EXISTS ?", clause.Table{Name: stmt.Table}).Error
			}); err != nil {
				return err
			}
		}

		return nil
	})
}

func (m Migrator) HasColumn(value interface{}, name string) bool {
	var count int
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(name); field != nil {
			name = field.DBName
		}

		if name != "" {
			if err := m.DB.Raw(
				"SELECT count(*) FROM sqlite_master WHERE type = ? AND tbl_name = ? AND (sql LIKE ? OR sql LIKE ? OR sql LIKE ?)",
				"table", stmt.Table, `%"`+name+`" %`, `%`+name+` %`, "%`"+name+"`%",
			).Row().Scan(&count); err != nil {
				return err
			}
		}
		return nil
	})
	return count > 0
}

func (m Migrator) AlterColumn(value interface{}, name string) error {
	return m.RunWithoutForeignKey(func() error {
		return m.RunWithValue(value, func(stmt *gorm.Statement) error {
			if field := stmt.Schema.LookUpField(name); field != nil {
				var (
					createSQL    string
					newTableName = stmt.Table + "__temp"
				)

				if err := m.DB.Raw("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "table", stmt.Table, stmt.Table).Row().Scan(&createSQL); err != nil {
					return err
				}

				reg, err := regexp.Compile("(`|'|\"| )" + field.DBName + "(`|'|\"| ) .*?,")
				if err == nil {
					tableReg, err := regexp.Compile(" ('|`|\"| )" + stmt.Table + "('|`|\"| ) ")
					if err != nil {
						return err
					}

					createSQL = tableReg.ReplaceAllString(createSQL, fmt.Sprintf(" `%v` ", newTableName))
					createSQL = reg.ReplaceAllString(createSQL, fmt.Sprintf("`%v` ?,", field.DBName))

					var columns []string
					columnTypes, _ := m.DB.Migrator().ColumnTypes(value)
					for _, columnType := range columnTypes {
						columns = append(columns, fmt.Sprintf("`%v`", columnType.Name()))
					}

					return m.DB.Transaction(func(tx *gorm.DB) error {
						queries := []string{
							createSQL,
							fmt.Sprintf("INSERT INTO `%v`(%v) SELECT %v FROM `%v`", newTableName, strings.Join(columns, ","), strings.Join(columns, ","), stmt.Table),
							fmt.Sprintf("DROP TABLE `%v`", stmt.Table),
							fmt.Sprintf("ALTER TABLE `%v` RENAME TO `%v`", newTableName, stmt.Table),
						}
						for _, query := range queries {
							if err := tx.Exec(query, m.FullDataTypeOf(field)).Error; err != nil {
								return err
							}
						}
						return nil
					})
				}
				return err
			}
			return fmt.Errorf("failed to alter field with name %v", name)
		})
	})
}

// ColumnTypes return columnTypes []gorm.ColumnType and execErr error
func (m Migrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	execErr := m.RunWithValue(value, func(stmt *gorm.Statement) (err error) {
		var (
			sqls   []string
			sqlDDL *ddl
		)

		if err := m.DB.Raw("SELECT sql FROM sqlite_master WHERE type IN ? AND tbl_name = ? AND sql IS NOT NULL order by type = ? desc", []string{"table", "index"}, stmt.Table, "table").Scan(&sqls).Error; err != nil {
			return err
		}

		if sqlDDL, err = parseDDL(sqls...); err != nil {
			return err
		}

		rows, err := m.DB.Session(&gorm.Session{}).Table(stmt.Table).Limit(1).Rows()
		if err != nil {
			return err
		}
		defer func() {
			err = rows.Close()
		}()

		var rawColumnTypes []*sql.ColumnType
		rawColumnTypes, err = rows.ColumnTypes()
		if err != nil {
			return err
		}

		for _, c := range rawColumnTypes {
			columnType := migrator.ColumnType{SQLColumnType: c}
			for _, column := range sqlDDL.columns {
				if column.NameValue.String == c.Name() {
					column.SQLColumnType = c
					columnType = column
					break
				}
			}
			columnTypes = append(columnTypes, columnType)
		}

		return err
	})

	return columnTypes, execErr
}

func (m Migrator) DropColumn(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if field := stmt.Schema.LookUpField(name); field != nil {
			name = field.DBName
		}

		var (
			createSQL    string
			newTableName = stmt.Table + "__temp"
		)

		if err := m.DB.Raw("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "table", stmt.Table, stmt.Table).Row().Scan(&createSQL); err != nil {
			return err
		}

		reg, err := regexp.Compile("(`|'|\"| )" + name + "(`|'|\"| ) .*?,")
		if err == nil {
			tableReg, err := regexp.Compile(" ('|`|\"| )" + stmt.Table + "('|`|\"| ) ")
			if err != nil {
				return err
			}

			createSQL = tableReg.ReplaceAllString(createSQL, fmt.Sprintf(" `%v` ", newTableName))
			createSQL = reg.ReplaceAllString(createSQL, "")

			var columns []string
			columnTypes, _ := m.DB.Migrator().ColumnTypes(value)
			for _, columnType := range columnTypes {
				if columnType.Name() != name {
					columns = append(columns, fmt.Sprintf("`%v`", columnType.Name()))
				}
			}

			return m.DB.Transaction(func(tx *gorm.DB) error {
				queries := []string{
					createSQL,
					fmt.Sprintf("INSERT INTO `%v`(%v) SELECT %v FROM `%v`", newTableName, strings.Join(columns, ","), strings.Join(columns, ","), stmt.Table),
					fmt.Sprintf("DROP TABLE `%v`", stmt.Table),
					fmt.Sprintf("ALTER TABLE `%v` RENAME TO `%v`", newTableName, stmt.Table),
				}
				for _, query := range queries {
					if err := tx.Exec(query).Error; err != nil {
						return err
					}
				}
				return nil
			})
		}
		return err
	})
}

func (m Migrator) CreateConstraint(interface{}, string) error {
	return ErrConstraintsNotImplemented
}

func (m Migrator) DropConstraint(interface{}, string) error {
	return ErrConstraintsNotImplemented
}

func (m Migrator) HasConstraint(value interface{}, name string) bool {
	var count int64
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		return m.DB.Raw(
			"SELECT count(*) FROM sqlite_master WHERE type = ? AND tbl_name = ? AND (sql LIKE ? OR sql LIKE ? OR sql LIKE ?)",
			"table", stmt.Table, `%CONSTRAINT "`+name+`" %`, `%CONSTRAINT `+name+` %`, "%CONSTRAINT `"+name+"`%",
		).Row().Scan(&count)
	})

	return count > 0
}

func (m Migrator) CurrentDatabase() (name string) {
	var null interface{}
	_ = m.DB.Raw("PRAGMA database_list").Row().Scan(&null, &name, &null)
	return
}

func (m Migrator) BuildIndexOptions(opts []schema.IndexOption, stmt *gorm.Statement) (results []interface{}) {
	for _, opt := range opts {
		str := stmt.Quote(opt.DBName)
		if opt.Expression != "" {
			str = opt.Expression
		}

		if opt.Collate != "" {
			str += " COLLATE " + opt.Collate
		}

		if opt.Sort != "" {
			str += " " + opt.Sort
		}
		results = append(results, clause.Expr{SQL: str})
	}
	return
}

func (m Migrator) CreateIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			opts := m.BuildIndexOptions(idx.Fields, stmt)
			values := []interface{}{clause.Column{Name: idx.Name}, clause.Table{Name: stmt.Table}, opts}

			createIndexSQL := "CREATE "
			if idx.Class != "" {
				createIndexSQL += idx.Class + " "
			}
			createIndexSQL += "INDEX ?"

			if idx.Type != "" {
				createIndexSQL += " USING " + idx.Type
			}
			createIndexSQL += " ON ??"

			if idx.Where != "" {
				createIndexSQL += " WHERE " + idx.Where
			}

			return m.DB.Exec(createIndexSQL, values...).Error
		}

		return fmt.Errorf("failed to create index with name %v", name)
	})
}

func (m Migrator) HasIndex(value interface{}, name string) bool {
	var count int
	_ = m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			name = idx.Name
		}

		if name != "" {
			if err := m.DB.Raw(
				"SELECT count(*) FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "index", stmt.Table, name,
			).Row().Scan(&count); err != nil {
				return err
			}
		}
		return nil
	})
	return count > 0
}

func (m Migrator) RenameIndex(value interface{}, oldName, newName string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		var sql string
		if err := m.DB.Raw("SELECT sql FROM sqlite_master WHERE type = ? AND tbl_name = ? AND name = ?", "index", stmt.Table, oldName).Row().Scan(&sql); err != nil {
			return err
		}
		if sql != "" {
			return m.DB.Exec(strings.Replace(sql, oldName, newName, 1)).Error
		}
		return fmt.Errorf("failed to find index with name %v", oldName)
	})
}

func (m Migrator) DropIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := stmt.Schema.LookIndex(name); idx != nil {
			name = idx.Name
		}

		return m.DB.Exec("DROP INDEX ?", clause.Column{Name: name}).Error
	})
}
