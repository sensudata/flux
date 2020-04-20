short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -3y  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "state", "location", "bottom_degrees", "surface_degrees"])

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
option now = () => 2020-02-22T18:00:00Z

short_hand_0
|> select_0()
