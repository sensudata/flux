## 2nd iteration of type inference

We need a 2nd iteration on type inference.
The goal of this 2nd iteration is twofold.

1. Provide a provably correct inference algorithm for Flux.
2. Represent a table stream as a list of record types and update the type signatures for each Flux transformation accordingly.

As outlined in [#1578](https://github.com/influxdata/flux/issues/1578), in order to be able to prove the correctness of type inference, we must break the mutual recursion between _type_ unification and _kind_ unification.
And to do that, we must replace our current record _kind_ with an extensible record _type_.

### Extensible Records

The Flux type system and the current Flux record _kind_ is based on the following [paper](https://caml.inria.fr/pub/papers/garrigue-structural_poly-fool02.pdf).
In this system record kinds have a set of properties along with their types.
They also contain two sets of labels.
One set specifies the properties that a record **must** have, the other specifies the properties it **may** have.
While this allows for polymorphic record access as well as polymorphic record literals, it lacks support for record extensions.
The following example cannot be typed in this system.

```
r = {a: 1, b: 2}
s = {r with c: 3} // s is a record with properties a, b, and c
```

In order to support examples like the one above, I propose we replace the current record kind with one that allows for extension such as the one in [this paper](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/scopedlabels.pdf).
Here the syntax of record types more closely aligns with the syntax of the values themselves.
In this system records are constructed via extension, starting from the empty record `{}`.

For example the record literal `{a:1, b:2}` is just shorthand for `{a:1 | {b:2 | {} }}` whose type is given by `{a:int | {b:int | {} }}`.

More generally, the polymorphic type of the extension operation `{_ with l:_}` is given by:

```
forall ['a, 'r] where 'r: Rec ('a, 'r) -> {l: 'a | 'r}
```

This type is read as, "for all type variables _'a_ and _'r_ where _'r_ is a record kind, we have a function type taking _'a_ and _'r_ as input and returning a record type extended from _'r_ by adding a field _l_ of type _'a_".
Selection and restriction operations are defined similarly, and update and rename can be expressed using these basic operations.

Unification of record types is particularly elegant and can be expressed using the following 3 rules:

1. `{a: 't | 'r} = {a: 'u | 's} => 't = 'u and 'r = 's`
2. `{a: 't | 'r} = {b: 'u | 's} => {a: 't | 'v} = 's and {b: 'u | 'v} = 'r`
3. `{a: 't | 'r} = {b: 'u | 'r} => error`

Note here `=` is synonymous for unify.

### Implications

The Flux type system extends Hindley-Milner with a _kind_ system.
Currently records are modeled as _kinds_ in this system.
Extensible records however fall entirely within the _type_ system, not the _kind_ system.
Describing records this way allows us to break the mutual recursion between type unification and kind unification.
Proving the correctness of type inference now follows directly from Hindley-Milner.

### Caveats

* Nulls

    Extensible records do not support accessing non-existent fields on an record.
    Such operations will result in a unification error.
    The Flux language specification states that accessing a non-existent field on a record will return an untyped null value.
    The spec will need to be updated accordingly as this will no longer be the case.

    Null values will still be possible at runtime but they must be typed nulls.
    The concept of an untyped null will need to be removed from the spec.

* Polymorphic Record Literals

    The current record _kinds_ allow for polymorphic record literals.
    For example, currently you can define a list of heterogeneous records like so:
    ```
    a = [ {a:1, b:"str"}, {b:"str", c:3.3} ]
    ```

    This would result in a unification error in the case of extensible records.
    What this means is that we lose some of our ability to represent polymorphic container types.

    This might be concerning since the 2nd goal of this iteration is to represent table streams as a list of **heterogeneous** record types.
    However such concern is unwarranted as streams cannot be constructed from atomic types the way list literals can.
    Instead they are obtained entirely via a side effect operation - a call to the database.
    This means that from the point of view of the type system, a table stream is just a list of **homogeneous** records.
    In other words, record fields or columns that are never accessed never need to be typed.

    Additionally if we assume transformations are given correct types, type inference can determine the columns and column types that the user expects from the database.
    Storage can then implicitly filter out any types that do not match the user's expectation.

### Typing flux functions using lists of records

The goal of this 2nd iteration is **not** to assign an accurate type to every Flux transformation.
As it stands, it is not possible to accurately describe certain Flux transformations with a type signature.
Instead we should give every transformation a type signature of the form:

```
forall ['r, 's] where 'r:Rec, 's:Rec (tables: ['r], ...) -> ['s]
```

In words, every Flux transformation is a function which takes in a list of record kinds and returns a list of record kinds.
This gives us parity with the current system while allowing us to iteratively assign more accurate and descriptive types, giving us the ability to catch more errors at compile time.

Below I attempt to assign accurate types to common Flux transformations.

* from

    ```
    forall ['r, 's] where 's: StrRec (org: str, bucket: str) -> [ {_time: time | _start: time | _stop: time | _value: 'r | 's} ]
    ```
    This polymorphic type signature says that `from` is a function that returns a list of records.
    These records have fields `_time`, `_start`, `_stop`, and `_value`.
    The `_value` field can be of any type.
    All other fields represented by the row variable `'s` must be of type _string_.
    Note the last part corresponds to the tag values of a series.

    An example of a program that will fail to compile is:
    ```
    from(bucket: "telegraf")
        |> range(start: -1h)
        |> filter(fn: (r) => r.host == 1) // error: {host: int | 't} is not of kind StrRec
    ```
    The reason for failed compilation will be clear soon, but basically `range` and `filter` are just pass through transformations.
    They don't extend or modify the types of their inputs.
    As a result `filter` expects to receive the same list type that `from` returns and so it expects `r.host` to be of type _string_.
    This contradicts the usage of `r.host` as an _int_.

* range

    ```
    forall ['r] where 'r:Rec (start: time, stop: time, tables: [ {_start: time | _stop: time | 'r} ]) -> [ {_start: time | _stop: time | 'r} ]
    ```
    Let's focus on the `tables` parameter.
    `range` takes a `tables` parameter that is a list of records that must have fields `_start` and `_stop` and returns the same list type.

* filter

    ```
    forall ['r] where 'r:Rec (fn: (r: 'r) -> bool, tables: ['r]) -> ['r]
    ```
    `filter` takes a list of records and returns that same list of records.

    From the type signatures of `range` and `filter` it should be clear why the above example program fails unification.

* map

    ```
    forall ['r, 's] where 'r:Rec, 's:Rec (fn: ('r) => 's, tables: [ 'r ]) -> [ 's ]
    ```

* keep

    `keep` and `drop` operations can be implemented using `map`:
    ```
    map(fn: (r) => ( {a: r.a | b: r.b | {} } )) // keep columns a and b
    map(fn: (r) => ( {r - a, b} ))              // drop columns a and b
    ```
    Assigning a type signature to the actual `keep` and `drop` functions unfortunatly is not so simple.
    These functions take the names of the columns they operate on via a `columns` parameter.
    Assigning an accurate type signature to these functions requires the notion of compile time constants or label types [#1122](https://github.com/influxdata/flux/issues/1122).

    Another option is a preprocessing step before type inference that rewrites the Flux AST by replacing calls to `keep` and `drop` with their `map` equivalents.

* aggregates

    Assigning accurate types to aggregates is also a bit challenging.
    This is due to the fact the aggregates will implicitly drop columns that are not part of the group key.
    However if aggregates didn't drop columns dynamically, we could type them like so:
    ```
    forall ['r] (column: str, tables: [ {_time: time | 'r} ]) -> [ 'r ]
    ```
    In words, aggregates take a list of records with a `_time` fields and return that same list minus the `_time` field.

Note that functions which remove columns from tables dynamically pose challenges to type inference.
A generic type can still be given to the function but any information we have on the input, we lose on the output.
`join` has the same problem, but for a different reason.
Conceptually `join` just concatenates records.
Unfortunately record concatenation is not supported by the type system and as a result all we can say about join is that it produces a generic list of records.
Hence any information we had about the records going into join, we lose going out.

#### What about group?

Group presents a different challenge.
While we can think of streams logically as a list of heterogeneous records, internally they are represented using a much more efficient columnar implementation.
So while the type signature of `group` is just:
```
forall ['r] where 'r:Rec (columns: [str], mode: str, tables: ['r]) -> ['r]
```
grouping can potentially place data of two different types in the same column, which will result in a type error at runtime.

A potential way to avoid this would be to create a new _kind_ describing any record whose type is known for every possible key.
In other words, a record _kind_ that enforces schema homogeneity.
As a result a query like `from() |> range() |> group(columns: ["host"])` would fail unification as the type of the `_value` column produced by `from` is unknown.
This might be too restrictive though.

#### What about pivot?

The same strategy we employed for `group` we can also employ for `pivot`.
This would give us partial information on the type of pivot's output.
However this too might be too restrictive.
Do users want to be able to group and pivot heterogeneous data?

#### What about other data sources?

Unlike InfluxDB, sources like `csv` and `sql` cannot have such a strict type signature.
Instead sources like these need to be defined to return lists of generic records `forall ['r] (...) -> [ 'r ]`.
However based on usage, type inference will be able to determine a more specific type for `'r`.
Each source can then perform their own type checking to ensure the data on disk matches with the type of `'r` that is inferred.

### Next Steps

So how do we proceed from here?

1. Replace our current record kinds with extensible record types.
2. Simplify kind system by removing mutual recursion between type unification and kind unification.
3. Give every Flux transformation the generalized type signature `forall ['r, 's] (tables: ['r], ...) -> ['s]`.
4. Refactor how row functions are type checked to allow for null values at runtime.

Completing steps 1-4 will accomplish the goals of this 2nd iteration on the type inference system.
Afterwards we can iteratively work to assign more accurate type signatures to the current Flux transformations in order to catch even more errors at compile time.
