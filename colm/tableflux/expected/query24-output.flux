short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -1y  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "state"])

select_0 = (tables=<-) => {
	return tables
		|> window( every: 1h )
		|> drop( columns: ["_start", "_time"] )
		|> rename( columns: { _stop: "_time" } )
		|> distinct( column: "state" )
		|> rename( columns: { _value: "state" } )
}
short_hand_0
|> select_0()
