module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1, h2, h3, ol, ul, li, a, span, header, footer, i, img, input)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (href, class, id, disabled, attribute, src, placeholder)
import List exposing (foldl)
import Time
import Task
import Random
import Http
import Delay
import Json.Decode as Decode exposing (Decoder, map2, map3, map8, int, bool, float, nullable, string, field, list, decodeString)
import Json.Encode
import Json.Decode.Pipeline exposing (required, optional)
import Round exposing (round)
import String exposing (toList)
import Char exposing (toCode)
import Url exposing (percentEncode)


type alias Model =
    { error : String
    , content : Maybe Content
    , sorting : String
    , currentPage : Int
    , prefix : String
    , cart : Cart
    }

type Content
  = Products ProductList
  | Product ProductDetail
  | CartPage

type alias ProductDetail =
  { uuid:          String
  , title:         String
  , description:   String
  , longtext:      String
  , category:      String
  , smallImageURL: String
  , largeImageURL: String
  , price:         Maybe Float
  }

type alias PaginatedMeta =
  { totalItems: Int
  , totalPages: Int
  , currentPage: Int
  , itemsPerPage: Int
  }

type alias ProductOverviewList = List ProductOverview

type alias ProductList =
  { data: ProductOverviewList
  , meta: PaginatedMeta
  }
type alias ProductOverview =
  { uuid             : String
  , title            : String
  , price            : Maybe Float
  }

type alias Cart = List CartItem
type alias CartItem =
  { product: ProductDetail
  , quantity: Int
  , moreInStock: Bool
  , inStock: Bool
  }

type alias CartChange =
  { action: CartAction
  , uuid: String
  }

type CartAction
  = Add
  | Remove


init : () -> ( Model, Cmd Msg )
init _ =
    ( { error = ""
      , content = Nothing
      , sorting = "name"
      , currentPage = 0
      , prefix = ""
      , cart = []
      }
    , Cmd.batch [ fetchProducts "name" 0 "", fetchCart ]
    )


type Msg
    = LoadProducts
    | LoadProduct String
    | PreviousPage
    | NextPage
    | SortByUuid
    | SortByPrice
    | SortByName
    | FilterByName
    | AddToCart String
    | RemoveFromCart String
    | GotProducts (Result Http.Error ProductList)
    | GotProduct (Result Http.Error ProductDetail)
    | CartGotChanged (Result Http.Error Cart)
    | ShowCart
    | GoToPage String
    | SetFilterPrefix String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of

        LoadProducts ->
            ( model, fetchProducts model.sorting model.currentPage model.prefix)

        LoadProduct uuid ->
            ( model, fetchProduct uuid )

        PreviousPage ->
            ( model, fetchProducts model.sorting (model.currentPage - 1)  model.prefix )

        NextPage ->
            ( model, fetchProducts model.sorting (model.currentPage + 1)  model.prefix )

        SortByUuid ->
            ( { model | sorting = "uuid" }, fetchProducts "uuid" model.currentPage model.prefix )

        SortByPrice ->
            ( { model | sorting = "price" }, fetchProducts "price" model.currentPage model.prefix )

        SortByName ->
            ( { model | sorting = "name" }, fetchProducts "name" model.currentPage model.prefix )

        FilterByName ->
            ( { model | sorting = "prefix" }, fetchProducts "prefix" model.currentPage model.prefix )

        AddToCart uuid ->
            ( model, updateCart (CartChange Add uuid) )

        RemoveFromCart uuid ->
            ( model, updateCart (CartChange Remove uuid) )

        GotProducts result ->
            case result of
                Ok pp ->
                    ( { model | content = Just (Products pp), currentPage = pp.meta.currentPage, error = "" }, Cmd.none )
                Err e ->
                    ( { model | content = Nothing, error = toString e }, Cmd.none )

        GotProduct result ->
            case result of
                Ok p ->
                    ( { model | content = Just (Product p), error = "" }, Cmd.none )
                Err e ->
                    ( { model | content = Nothing, error = toString e }, Cmd.none )

        CartGotChanged result ->
            case result of
                Ok c ->
                    ( { model | cart = c, error = "" }, Cmd.none )
                Err e ->
                    ( { model | content = Nothing, error = toString e }, Cmd.none )

        ShowCart ->
            ( { model | content = Just CartPage }, Cmd.none  )

        GoToPage page ->
            ( model, fetchProducts model.sorting (Maybe.withDefault 0 (String.toInt page) - 1) model.prefix )

        SetFilterPrefix prefix ->
            ( { model | prefix = prefix }, fetchProducts model.sorting model.currentPage prefix )


toString : Http.Error -> String
toString error =
    case error of
        Http.BadUrl url ->
            url ++ " is bad"

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network Error"

        Http.BadStatus status ->
            String.fromInt status ++ " Status"

        Http.BadBody msg ->
            msg


view : Model -> Html Msg
view model =
    let
        (showingCart, showingCatalog) = case model.content of
            Just (Products pp) -> (False, True)
            Just (Product p) -> (False, False)
            Just CartPage -> (True, False)
            Nothing -> (False, False)
    in
    div [class "mdl-layout mdl-layout--fixed-header"]
    [ header [ class "mdl-layout__header mdl-layout__header--waterfall custom-header"]
      [ div [ class "mdl-layout__header-row custom-header-row"]
        [ span [ class "mdl-layout__title"] [ text "Event Thingy Store" ]
        , div [ class "mdl-layout-spacer"] []
        , div [] [ span [ class "mdl-badge custom-header-cart", attribute "data-badge" (itemsInCart model.cart), onClick ShowCart ] [ text "Cart" ] ]
        , button [ class "mdl-button mdl-button--raised mdl-button--accent", onClick ShowCart, disabled showingCart ] [ text "show Cart" ]
        , button [ class "mdl-button mdl-button--raised mdl-button--accent", onClick LoadProducts, disabled showingCatalog ] [ text "show products" ]
        ]
      ]
      , div [ class "custom-header-error"] [ text model.error ]
      , div [ class "mdl-layout__content", id "main"] [ renderContent model ]
    ]

itemsInCart : Cart -> String
itemsInCart cart =
    let
        inStock = itemsInCartInStock cart
        outOfStock = itemsInCartOutOfStock cart
    in
        if outOfStock == 0 then
            String.fromInt inStock
        else
            String.fromInt (inStock + outOfStock) ++ "-" ++ String.fromInt outOfStock

itemsInCartInStock : Cart -> Int
itemsInCartInStock cart =
    let
        ls = List.map (\e -> if e.inStock then 1 else 0 ) cart
    in
        foldl (+) 0 ls
        
itemsInCartOutOfStock : Cart -> Int
itemsInCartOutOfStock cart =
    let
        ls = List.map (\e -> if not e.inStock then 1 else 0 ) cart
    in
        foldl (+) 0 ls

renderContent : Model -> Html Msg
renderContent model =
    case model.content of
        Just (Products pp) -> renderProducts pp model.sorting
        Just (Product p) -> renderProductDetail p
        Just CartPage -> renderCart model.cart
        Nothing -> text "" 


renderProducts : ProductList -> String -> Html Msg
renderProducts lst sorting =
    let
        prevEnabled = lst.meta.currentPage > 0
        nextEnabled = lst.meta.currentPage < (lst.meta.totalPages - 1)
        uuidDisabled = sorting == "uuid"
        priceDisabled = sorting == "price"
        nameDisabled = sorting == "name" || sorting == "prefix"
        prefixDisabled = sorting == "prefix"
        pagesText = if lst.meta.totalPages == 0
                    then text " No luck! "
                    else text (" Page " ++ String.fromInt (lst.meta.currentPage + 1) ++ " from " ++ String.fromInt(lst.meta.totalPages))
    in
    div [ class "mdl-grid" ] [
        div [class "mdl-cell mdl-cell--12-col"]
        [ button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick PreviousPage, disabled (not prevEnabled) ] [ i [ class "material-icons" ] [ text "remove" ] ]
        , span [ class "mdl-title custom-page-number"]
          [ pagesText
          , input [ class "custom-go-to-page", placeholder "Go to page", onInput GoToPage ] []
          ]
        , button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick NextPage, disabled (not nextEnabled) ] [ i [ class "material-icons" ] [ text "add" ] ]
        , span [ class "custom-sorting" ]
          [ text "Sort by "
          , sortProductsButton SortByName nameDisabled "Name"
          , sortProductsButton SortByUuid uuidDisabled "Uuid"
          , sortProductsButton SortByPrice priceDisabled "Price"
          ]
        , span [ class "custom-sorting" ]
          [ text "Filter by "
          , sortProductsButton FilterByName prefixDisabled "Prefix"
          , input [ onInput SetFilterPrefix, disabled (not prefixDisabled) ] []
          ]
        ]
        , ul [ class "product-list mdl-list" ] (List.map (\l -> renderProduct l ) lst.data )
    ]


sortProductsButton : Msg -> Bool -> String -> Html Msg
sortProductsButton click grey label =
    button
        [ class "mdl-button mdl-js-button mdl-button--raised mdl-button--colored", onClick click, disabled grey ]
        [ text label ]


showProductButton : String -> Html Msg
showProductButton uuid =
    button
      [ class "mdl-button mdl-js-button mdl-button--raised mdl-button--colored", onClick (LoadProduct uuid) ]
      [ text "show Details" ]


renderProduct : ProductOverview  -> Html Msg
renderProduct product =
    li [ class "mdl-list__item mdl-list__item--two-line" ]
    [ span [ class "mdl-list__item-primary-content" ]
      [ img [ class "custom-list-image",src (productImage product.uuid 100 50) ] []
      , span [ onClick (LoadProduct product.uuid) ] [ text product.title ]
      , span [ class "mdl-list__item-sub-title" ] [ text ("price: " ++ formatPrice product.price) ]
      ]
      , span [ class "mdl-list__item-secondary-content" ] [ showProductButton product.uuid ]
      , span [ class "mdl-list__item-secondary-content" ] [ addToCartButton product.uuid ]
    ]    


productImage : String -> Int -> Int -> String
productImage uuid width height =
    "https://picsum.photos/id/" ++ String.fromInt (modBy 50 (reduceUuid uuid)) ++ "/" ++ String.fromInt width ++ "/" ++ String.fromInt height


reduceUuid : String -> Int
reduceUuid uuid =
    List.foldl (\x a -> x + a) 0 (List.map toCode (toList uuid))


formatPrice : Maybe Float -> String
formatPrice price =
    round 2 (Maybe.withDefault 0.0 price) ++ "€"


renderProductDetail : ProductDetail  -> Html Msg
renderProductDetail product =
    div [ class "mdl-grid"]
    [ div [ class "mdl-cell mdl-cell--12-col" ] [ h1 [] [ text product.title ] ]
    , div [ class "mdl-cell mdl-cell--8-col" ] [ img [ class "custom-detail-image",src (productImage product.uuid 400 200)] [] ]
    , div [ class "mdl-cell mdl-cell--4-col" ] 
      [ span [ class "custom-detail-block" ] [ text ("id: " ++ product.uuid) ]
      , span [ class "mdl-typography--headline custom-detail-block" ] [ text product.description ]
      , span [ class "custom-detail-block" ] [ text product.longtext ]
      , span [ class "mdl-typography--display-1 custom-detail-block" ] [ text (formatPrice product.price) ] 
      , span [] [ addToCartButton product.uuid ] 
      ]
    ]


addToCartButton : String -> Html Msg
addToCartButton uuid =
    button
      [ class "mdl-button mdl-js-button mdl-button--raised mdl-js-ripple-effect mdl-button--accent", onClick (AddToCart uuid) ]
      [ text "add to cart" ]


renderCart : Cart -> Html Msg
renderCart cart =
    ul [ class "product-list mdl-list" ] (List.map (\l -> renderCartItem l ) cart )


renderCartItem : CartItem -> Html Msg
renderCartItem item =
    let
        uuid = item.product.uuid
        product = item.product
        quantity = item.quantity
        inStock = item.inStock
        outOfStock = not item.inStock
        moreInStock = item.moreInStock
        stockText = if not item.inStock then text "Out of Stock"
                    else if not item.moreInStock then text "All stock in cart"
                    else text ""
    in
        li [ class "mdl-list__item mdl-list__item--two-line" ]
        [ span [ class "mdl-list__item-primary-content" ]
          [ img [ class "custom-list-image",src (productImage product.uuid 100 50) ] []
          , span [ onClick (LoadProduct product.uuid) ] [ text product.title ]
          , span [ class "mdl-list__item-sub-title" ] [ text ("price: " ++ formatPrice product.price) ]
          ]
          , span []
            [ stockText
            , button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick (RemoveFromCart uuid) ] [ i [ class "material-icons" ] [ text "remove" ] ]
            , span [] [ text (String.fromInt quantity) ]
            , button [ class "mdl-button mdl-js-button mdl-button--fab mdl-button--colored", onClick (AddToCart uuid), disabled (not moreInStock) ] [ i [ class "material-icons" ] [ text "add" ] ]
            ]
        ]    

fetchProducts : String -> Int -> String -> Cmd Msg
fetchProducts sor currentPage pre =
    let
        sorting = case sor of
            "prefix" -> "name" 
            _ -> sor
        prefix = percentEncode pre
    in    
    Http.get
        { url = "http://localhost:8080/products?sort=" ++ sorting ++ "&prefix=" ++ prefix ++ "&page=" ++ String.fromInt currentPage
        , expect = Http.expectJson GotProducts productListDecoder
        }


fetchProduct : String -> Cmd Msg
fetchProduct uuid =
    Http.get
        { url = "http://localhost:8080/product?uuid=" ++ uuid
        , expect = Http.expectJson GotProduct productDetailDecoder
        }


updateCart : CartChange -> Cmd Msg
updateCart cartChange =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = "http://localhost:8080/cart"
        , body = Http.jsonBody (encodeCartChange cartChange)
        , expect = Http.expectJson CartGotChanged cartDecoder
        , timeout = Nothing
        , tracker = Nothing
        }


fetchCart : Cmd Msg
fetchCart =
    Http.riskyRequest
        { method = "GET"
        , headers = []
        , url = "http://localhost:8080/cart"
        , body = Http.emptyBody
        , expect = Http.expectJson CartGotChanged cartDecoder
        , timeout = Nothing
        , tracker = Nothing
        }


encodeCartChange : CartChange -> Json.Encode.Value
encodeCartChange cc  =
    Json.Encode.object
        [ ( "action", Json.Encode.int (cartActionToString cc.action) )
        , ( "uuid", Json.Encode.string cc.uuid )
        ]


cartActionToString : CartAction -> Int
cartActionToString c =
    case c of
-- needs to match the values from protobuf
      Add -> 0
      Remove -> 1


subscriptions : Model -> Sub Msg
subscriptions _ = Sub.batch []


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }

paginatedMetaDecoder : Decoder PaginatedMeta
paginatedMetaDecoder =
    Decode.succeed PaginatedMeta
        |> required "total_items" int
        |> required "total_pages" int
        |> required "current_page" int
        |> required "items_per_page" int

productListDecoder : Decoder ProductList
productListDecoder = 
    Decode.succeed ProductList
        |> required "data" (list productDecoder)
        |> required "meta" paginatedMetaDecoder

productDecoder : Decoder ProductOverview
productDecoder =
    Decode.succeed ProductOverview
        |> required "uuid" string
        |> required "title" string
        |> optional "price" (nullable float) Nothing


productDetailDecoder : Decoder ProductDetail
productDetailDecoder = 
    Decode.succeed ProductDetail
        |> required "uuid" string
        |> required "title" string
        |> required "description" string
        |> required "longtext" string
        |> required "smallImageURL" string
        |> required "largeImageURL" string
        |> required "description" string
        |> optional "price" (nullable float) Nothing

cartDecoder : Decoder Cart
cartDecoder = list cartItemDecoder


cartItemDecoder : Decoder CartItem
cartItemDecoder =
    Decode.succeed CartItem
        |> required "product" productDetailDecoder
        |> required "quantity" int
        |> required "moreInStock" bool
        |> required "inStock" bool
