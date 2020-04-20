short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -1y  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "state", "location", "bottom_degrees"])

select_0 = (tables=<-) => {
	return tables
		|> group( columns: ["state"] )
		|> bottom( columns: ["bottom_degrees"], n: 2 )
}
short_hand_0
|> select_0()
