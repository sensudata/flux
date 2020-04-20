short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -1y  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "location", "surface_degrees"])

select_0 = (tables=<-) => {
	return tables
		|> top( columns: ["surface_degrees"], n: 5 )
}
short_hand_0
|> select_0()
