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
aggregate_0 = (tables=<-) => {
	grouping = tables
		|> window( every: 1h )
	local_1 = grouping
		|> keep( columns: ["_stop", "state"] )
		|> rename( columns: {state: "count_state"} )
		|> count( column: "count_state" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
	return local_1
		|> drop(columns: ["__id"] )
		|> rename( columns: { _stop: "_time" } )
}

short_hand_0
|> select_0()
|> timeShift( duration: -1m, columns: ["_time"] )
|> aggregate_0()
