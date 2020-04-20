short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -1y  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "location", "bottom_degrees", "surface_degrees"])

select_0 = (tables=<-) => {
	return tables
		|> map( fn: (r) => ({r with _value: 0}) )
		|> group( columns: ["location"] )
		|> window( every: 1h )
		|> drop( columns: ["_start", "_time"] )
		|> rename( columns: { _stop: "_time" } )
		|> first()
		|> drop( columns: ["_value"] )
}
short_hand_0
|> select_0()
