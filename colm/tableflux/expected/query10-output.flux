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

option now = () => 2020-02-22T18:00:00Z

short_hand_0
