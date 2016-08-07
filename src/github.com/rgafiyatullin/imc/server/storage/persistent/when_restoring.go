package persistent

import (
	"github.com/rgafiyatullin/imc/server/storage/inmemory/bucket/data"
)

func (this *state) whenRestoring() {
	this.fetchAndRestoreValues()

	this.status = stSaving
	this.chans.restore <- NewRestoreComplete()
	this.db.Close()
	this.startWriter()
}

func (this *state) fetchAndRestoreValues() {
	rows, err := this.db.Query("SELECT k, v, t, e FROM keys")
	if err != nil {
		this.actorCtx.Log().Error("Error restoring string values: %v", err)
	} else {
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				this.actorCtx.Log().Error("string> error: %v", err)
			} else if len(columns) != 4 {
				this.actorCtx.Log().Error("string> wrong columns count: %d", len(columns))
			} else {
				var key string
				var keyType int
				var val []byte
				var extId int
				rows.Scan(&key, &val, &keyType, &extId)
				switch keyType {
				case ktString:
					this.actorCtx.Log().Debug("string> '%v' -> '%v'", key, val)
					this.chans.restore <- NewRestoreString(key, data.NewScalar(val))
				case ktList:
					this.actorCtx.Log().Debug("list> '%v' -> %v[...]", key, extId)
					listRows, err := this.db.Query("SELECT value FROM lists WHERE list_id = ? ORDER BY idx ASC", extId)
					if err == nil {
						l := data.NewList()
						for listRows.Next() {
							var v []byte
							listRows.Scan(&v)
							l.PushBack(v)
						}

						if l.Len() > 0 {
							this.chans.restore <- NewRestoreList(key, l)
						}
					} else {
						this.actorCtx.Log().Warning("list> '%v' -> Err(%v)", key, err)
					}

				case ktMap:
					this.actorCtx.Log().Debug("map> '%v' -> %v{...}", key, extId)
					mapPairs, err := this.db.Query("SELECT k, value FROM maps WHERE map_id = ?", extId)
					if err == nil {
						d := data.NewDict()

						for mapPairs.Next() {
							var key string
							var val []byte

							mapPairs.Scan(&key, &val)
							d.Set(key, val)
						}

						if d.Len() > 0 {
							this.chans.restore <- NewRestoreDict(key, d)
						}
					} else {
						this.actorCtx.Log().Warning("map> '%v' -> Err(%v)", key, err)
					}
				}

			}
		}
	}
}
