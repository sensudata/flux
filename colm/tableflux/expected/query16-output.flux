short_hand_0 = from( bucket: "tableflux" )
	|> range( start: -3h  )
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "state", "location", "bottom_degrees", "surface_degrees"])

aggregate_0 = (tables=<-) => {
	grouping = tables
		|> group( columns: ["state"] )
	local_1 = grouping
		|> keep( columns: ["bottom_degrees", "state"] )
		|> rename( columns: {bottom_degrees: "min_bottom_degrees"} )
		|> min( column: "min_bottom_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
	local_2 = grouping
		|> keep( columns: ["bottom_degrees", "state"] )
		|> rename( columns: {bottom_degrees: "max_bottom_degrees"} )
		|> max( column: "max_bottom_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
		|> drop( columns: ["bottom_degrees", "state"]  )
	local_3 = grouping
		|> keep( columns: ["surface_degrees", "state"] )
		|> rename( columns: {surface_degrees: "min_surface_degrees"} )
		|> min( column: "min_surface_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
		|> drop( columns: ["surface_degrees", "state"]  )
	local_4 = grouping
		|> keep( columns: ["surface_degrees", "state"] )
		|> rename( columns: {surface_degrees: "max_surface_degrees"} )
		|> max( column: "max_surface_degrees" )
		|> map( fn: (r) => ( {r with __id: 1} ) )
		|> group()
		|> cumulativeSum( columns: ["__id"] )
		|> drop( columns: ["surface_degrees", "state"]  )
	local_5 = join(
		tables: {
			local_1: local_1,
			local_2: local_2
		},
		on: ["__id"]
	)
	local_6 = join(
		tables: {
			local_5: local_5,
			local_3: local_3
		},
		on: ["__id"]
	)
	local_7 = join(
		tables: {
			local_6: local_6,
			local_4: local_4
		},
		on: ["__id"]
	)
	return local_7
		|> drop(columns: ["__id"] )
}

option now = () => 2020-02-22T18:00:00Z

short_hand_0
|> aggregate_0()
