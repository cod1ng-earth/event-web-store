module Products exposing (..)


import Json.Decode as JD
import Json.Encode as JE


(<$>) : (a -> b) -> JD.Decoder a -> JD.Decoder b
(<$>) =
    JD.map


(<*>) : JD.Decoder (a -> b) -> JD.Decoder a -> JD.Decoder b
(<*>) f v =
    f |> JD.andThen (\x -> x <$> v)


optionalDecoder : JD.Decoder a -> JD.Decoder (Maybe a)
optionalDecoder decoder =
    JD.oneOf
        [ JD.map Just decoder
        , JD.succeed Nothing
        ]


requiredFieldDecoder : String -> a -> JD.Decoder a -> JD.Decoder a
requiredFieldDecoder name default decoder =
    withDefault default (JD.field name decoder)


optionalFieldDecoder : String -> JD.Decoder a -> JD.Decoder (Maybe a)
optionalFieldDecoder name decoder =
    optionalDecoder (JD.field name decoder)


repeatedFieldDecoder : String -> JD.Decoder a -> JD.Decoder (List a)
repeatedFieldDecoder name decoder =
    withDefault [] (JD.field name (JD.list decoder))


withDefault : a -> JD.Decoder a -> JD.Decoder a
withDefault default decoder =
    JD.oneOf
        [ decoder
        , JD.succeed default
        ]


optionalEncoder : String -> (a -> JE.Value) -> Maybe a -> Maybe (String, JE.Value)
optionalEncoder name encoder v =
    case v of
        Just x ->
            Just ( name, encoder x )

        Nothing ->
            Nothing


requiredFieldEncoder : String -> (a -> JE.Value) -> a -> a -> Maybe ( String, JE.Value )
requiredFieldEncoder name encoder default v =
    if v == default then
        Nothing
    else
        Just ( name, encoder v )


repeatedFieldEncoder : String -> (a -> JE.Value) -> List a -> Maybe (String, JE.Value)
repeatedFieldEncoder name encoder v =
    case v of
        [] ->
            Nothing
        _ ->
            Just (name, JE.list <| List.map encoder v)


type CartChangeAction
    = Add -- 0
    | Remove -- 1


cartChangeActionDecoder : JD.Decoder CartChangeAction
cartChangeActionDecoder =
    let
        lookup s =
            case s of
                "add" ->
                    Add

                "remove" ->
                    Remove

                _ ->
                    Add
    in
        JD.map lookup JD.string


cartChangeActionDefault : CartChangeAction
cartChangeActionDefault = Add


cartChangeActionEncoder : CartChangeAction -> JE.Value
cartChangeActionEncoder v =
    let
        lookup s =
            case s of
                Add ->
                    "add"

                Remove ->
                    "remove"

    in
        JE.string <| lookup v


type alias CatalogPage =
    { data : List Pb.Product -- 1
    , meta : Maybe Pb.CatalogPageMetadata -- 2
    }


catalogPageDecoder : JD.Decoder CatalogPage
catalogPageDecoder =
    JD.lazy <| \_ -> CatalogPage
        <$> (repeatedFieldDecoder "data" pb_ProductDecoder)
        <*> (optionalFieldDecoder "meta" pb_CatalogPageMetadataDecoder)


catalogPageEncoder : CatalogPage -> JE.Value
catalogPageEncoder v =
    JE.object <| List.filterMap identity <|
        [ (repeatedFieldEncoder "data" pb_ProductEncoder v.data)
        , (optionalEncoder "meta" pb_CatalogPageMetadataEncoder v.meta)
        ]


type alias CatalogPageMetadata =
    { totalItems : Int -- 1
    , totalPages : Int -- 2
    , currentPage : Int -- 3
    , setPageTo : Int -- 4
    , sorting : String -- 5
    , filtering : String -- 6
    , itemsPerPage : Int -- 7
    }


catalogPageMetadataDecoder : JD.Decoder CatalogPageMetadata
catalogPageMetadataDecoder =
    JD.lazy <| \_ -> CatalogPageMetadata
        <$> (requiredFieldDecoder "totalItems" 0 JD.int)
        <*> (requiredFieldDecoder "totalPages" 0 JD.int)
        <*> (requiredFieldDecoder "currentPage" 0 JD.int)
        <*> (requiredFieldDecoder "setPageTo" 0 JD.int)
        <*> (requiredFieldDecoder "sorting" "" JD.string)
        <*> (requiredFieldDecoder "filtering" "" JD.string)
        <*> (requiredFieldDecoder "itemsPerPage" 0 JD.int)


catalogPageMetadataEncoder : CatalogPageMetadata -> JE.Value
catalogPageMetadataEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "totalItems" JE.int 0 v.totalItems)
        , (requiredFieldEncoder "totalPages" JE.int 0 v.totalPages)
        , (requiredFieldEncoder "currentPage" JE.int 0 v.currentPage)
        , (requiredFieldEncoder "setPageTo" JE.int 0 v.setPageTo)
        , (requiredFieldEncoder "sorting" JE.string "" v.sorting)
        , (requiredFieldEncoder "filtering" JE.string "" v.filtering)
        , (requiredFieldEncoder "itemsPerPage" JE.int 0 v.itemsPerPage)
        ]


type alias Product =
    { uuid : String -- 1
    , title : String -- 2
    , description : String -- 3
    , longtext : String -- 4
    , category : String -- 8
    , smallImageURL : String -- 5
    , largeImageURL : String -- 6
    , price : Float -- 7
    }


productDecoder : JD.Decoder Product
productDecoder =
    JD.lazy <| \_ -> Product
        <$> (requiredFieldDecoder "uuid" "" JD.string)
        <*> (requiredFieldDecoder "title" "" JD.string)
        <*> (requiredFieldDecoder "description" "" JD.string)
        <*> (requiredFieldDecoder "longtext" "" JD.string)
        <*> (requiredFieldDecoder "category" "" JD.string)
        <*> (requiredFieldDecoder "smallImageURL" "" JD.string)
        <*> (requiredFieldDecoder "largeImageURL" "" JD.string)
        <*> (requiredFieldDecoder "price" 0.0 JD.float)


productEncoder : Product -> JE.Value
productEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "uuid" JE.string "" v.uuid)
        , (requiredFieldEncoder "title" JE.string "" v.title)
        , (requiredFieldEncoder "description" JE.string "" v.description)
        , (requiredFieldEncoder "longtext" JE.string "" v.longtext)
        , (requiredFieldEncoder "category" JE.string "" v.category)
        , (requiredFieldEncoder "smallImageURL" JE.string "" v.smallImageURL)
        , (requiredFieldEncoder "largeImageURL" JE.string "" v.largeImageURL)
        , (requiredFieldEncoder "price" JE.float 0.0 v.price)
        ]


type alias ProductUpdate =
    { old : Maybe Pb.Product -- 1
    , new : Maybe Pb.Product -- 2
    }


productUpdateDecoder : JD.Decoder ProductUpdate
productUpdateDecoder =
    JD.lazy <| \_ -> ProductUpdate
        <$> (optionalFieldDecoder "old" pb_ProductDecoder)
        <*> (optionalFieldDecoder "new" pb_ProductDecoder)


productUpdateEncoder : ProductUpdate -> JE.Value
productUpdateEncoder v =
    JE.object <| List.filterMap identity <|
        [ (optionalEncoder "old" pb_ProductEncoder v.old)
        , (optionalEncoder "new" pb_ProductEncoder v.new)
        ]


type alias CheckoutContext =
    { checkoutContext : CheckoutContext
    }


type CheckoutContext
    = CheckoutContextUnspecified
    | CartChange Pb.CartChange
    | Stock Pb.Stock
    | ProductUpdate Pb.ProductUpdate
    | CartOrder Pb.CartOrder


checkoutContextDecoder : JD.Decoder CheckoutContext
checkoutContextDecoder =
    JD.lazy <| \_ -> JD.oneOf
        [ JD.map CartChange (JD.field "cartChange" pb_CartChangeDecoder)
        , JD.map Stock (JD.field "stock" pb_StockDecoder)
        , JD.map ProductUpdate (JD.field "productUpdate" pb_ProductUpdateDecoder)
        , JD.map CartOrder (JD.field "cartOrder" pb_CartOrderDecoder)
        , JD.succeed CheckoutContextUnspecified
        ]


checkoutContextEncoder : CheckoutContext -> Maybe ( String, JE.Value )
checkoutContextEncoder v =
    case v of
        CheckoutContextUnspecified -> Nothing
        CartChange x ->
            Just ( "cartChange", pb_CartChangeEncoder x )
        Stock x ->
            Just ( "stock", pb_StockEncoder x )
        ProductUpdate x ->
            Just ( "productUpdate", pb_ProductUpdateEncoder x )
        CartOrder x ->
            Just ( "cartOrder", pb_CartOrderEncoder x )


checkoutContextDecoder : JD.Decoder CheckoutContext
checkoutContextDecoder =
    JD.lazy <| \_ -> CheckoutContext
        <$> checkoutContextDecoder


checkoutContextEncoder : CheckoutContext -> JE.Value
checkoutContextEncoder v =
    JE.object <| List.filterMap identity <|
        [ (checkoutContextEncoder v.checkoutContext)
        ]


type alias CartChange =
    { cartID : String -- 3
    , uuid : String -- 1
    , action : Pb.CartChangeAction -- 2
    }


cartChangeDecoder : JD.Decoder CartChange
cartChangeDecoder =
    JD.lazy <| \_ -> CartChange
        <$> (requiredFieldDecoder "cartID" "" JD.string)
        <*> (requiredFieldDecoder "uuid" "" JD.string)
        <*> (requiredFieldDecoder "action" pb_CartChangeActionDefault pb_CartChangeActionDecoder)


cartChangeEncoder : CartChange -> JE.Value
cartChangeEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "cartID" JE.string "" v.cartID)
        , (requiredFieldEncoder "uuid" JE.string "" v.uuid)
        , (requiredFieldEncoder "action" pb_CartChangeActionEncoder pb_CartChangeActionDefault v.action)
        ]


type alias CartOrder =
    { cartID : String -- 1
    }


cartOrderDecoder : JD.Decoder CartOrder
cartOrderDecoder =
    JD.lazy <| \_ -> CartOrder
        <$> (requiredFieldDecoder "cartID" "" JD.string)


cartOrderEncoder : CartOrder -> JE.Value
cartOrderEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "cartID" JE.string "" v.cartID)
        ]


type alias Stock =
    { uuid : String -- 1
    , quantity : Int -- 2
    }


stockDecoder : JD.Decoder Stock
stockDecoder =
    JD.lazy <| \_ -> Stock
        <$> (requiredFieldDecoder "uuid" "" JD.string)
        <*> (requiredFieldDecoder "quantity" 0 JD.int)


stockEncoder : Stock -> JE.Value
stockEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "uuid" JE.string "" v.uuid)
        , (requiredFieldEncoder "quantity" JE.int 0 v.quantity)
        ]
