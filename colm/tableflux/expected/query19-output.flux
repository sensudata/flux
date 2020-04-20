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

_add_group_row_ids = (tables=<-) => {
	return tables
		|> map(fn: (r) => ({ r with row_id: 1}))
		|> cumulativeSum( columns: ["row_id"] )
		|> group()
		|> map(fn: (r) => ({ r with group_id: 1}))
		|> cumulativeSum( columns: ["group_id"] )
		|> map(fn: (r) => ({r with
			group_id: r.group_id - r.row_id}))
		|> difference( columns: ["group_id"], keepFirst: true )
		|> map(fn: (r) => ({r with
			group_id: if r.group_id > 0 then 1 else 0 }))
		|> cumulativeSum( columns: ["group_id"] )
}
select_0 = (tables=<-) => {
	with_ids = tables
		|> group( columns: ["state"] )
		|> window( every: 1h )
		|> drop( columns: ["_start", "_time"] )
		|> rename( columns: { _stop: "_time" } )
		|> _add_group_row_ids()
	grouped_fn = with_ids
		|> group( columns: ["state"] )
		|> window( every: 1h )
		|> drop( columns: ["_start", "_time"] )
		|> rename( columns: { _stop: "_time" } )
		|> min( column: "bottom_degrees" )
	fn_values = grouped_fn 
		|> group()
		|> sort( columns: ["group_id"] )
		|> tableFind( fn: (key) => (true) ) 
		|> getColumn( column: "bottom_degrees" )
	return with_ids
		|> filter( fn: (r) =>
			( r.bottom_degrees == fn_values[r.group_id] ) )
}
short_hand_0
|> select_0()
