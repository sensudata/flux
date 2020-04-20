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
		|> rename( columns: {bottom_degrees: "min_bottom_degrees"} )
		|> min( column: "min_bottom_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
	local_2 = grouping
		|> keep( columns: ["bottom_degrees"] )
		|> rename( columns: {bottom_degrees: "max_bottom_degrees"} )
		|> max( column: "max_bottom_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
		|> drop( columns: ["bottom_degrees"]  )
	local_3 = grouping
		|> keep( columns: ["bottom_degrees"] )
		|> rename( columns: {bottom_degrees: "mean_bottom_degrees"} )
		|> mean( column: "mean_bottom_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
		|> drop( columns: ["bottom_degrees"]  )
	local_4 = join(
		tables: {
			local_1: local_1,
			local_2: local_2
		},
		on: ["__id"]
	)
	local_5 = join(
		tables: {
			local_4: local_4,
			local_3: local_3
		},
		on: ["__id"]
	)
	return local_5
		|> drop(columns: ["__id"] )
}

option now = () => 2020-02-22T18:00:00Z

short_hand_0
|> aggregate_0()
