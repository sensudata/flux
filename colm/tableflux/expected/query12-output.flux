short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -3h  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "bottom_degrees"])

aggregate_0 = (tables=<-) => {
	grouping = tables
	local_1 = grouping
		|> keep( columns: ["bottom_degrees"] )
		|> rename( columns: {bottom_degrees: "count_bottom_degrees"} )
		|> count( column: "count_bottom_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
	return local_1
		|> drop(columns: ["__id"] )
}

option now = () => 2020-02-22T18:00:00Z

short_hand_0
|> aggregate_0()
