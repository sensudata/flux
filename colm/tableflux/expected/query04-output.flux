short_hand_0 = from( bucket: "tableflux" )
	|> range( start: 0 )
	|> last()
	|> filter( fn: (r) => ( r._measurement == "h2o_temperature") )
	|> pivot(
		rowKey:["_time"],
		columnKey: ["_field"],
		valueColumn: "_value"
	)
	|> group( columns: [] )
	|> drop( columns: ["_start", "_stop", "_measurement"] )

short_hand_0
