short_hand_0 = from( bucket: "tableflux" )
	|> range( start: 0 )
	|> last()
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> filter( fn: (r) => r.location == "coyote_creek" or r.location == "puget_sound" )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> keep( columns: ["_time", "time", "location", "surface_degrees"])

short_hand_0
