module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1, h2, h3, ol, ul, li, a, span, header, footer)
import Html.Events exposing (onClick)
import Html.Attributes exposing (href, class, id, disabled)
import List exposing (foldl)
import Time
import Task
import Random
import Http
import Delay
import Json.Decode as Decode exposing (Decoder, map2, map3, map8, int, float, nullable, string, field, list, decodeString)
import Json.Encode
import Json.Decode.Pipeline exposing (required, optional)


type alias Model =
    { error : String
    , content : Maybe Content
    , pageNumber: Int
    , cart : Cart
    }

type Content
  = Products Products
  | Product ProductDetail

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

type alias Products = List ProductOverview
type alias ProductOverview =
  { uuid             : String
  , title            : String
  , price            : Maybe Float
  }

type alias Cart = List CartItem
type alias CartItem =
  { product: ProductDetail
  , quantiy: Int
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
      , pageNumber = 0
      , cart = []
      }
    , Cmd.batch [ fetchProducts 0 ]
    )


type Msg
    = LoadProducts
    | LoadProduct String
    | PreviousPage
    | NextPage
    | AddToCart String
    | GotProducts (Result Http.Error Products)
    | GotProduct (Result Http.Error ProductDetail)
    | AddedToCart (Result Http.Error Cart)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of

        LoadProducts ->
            ( model, fetchProducts model.pageNumber)

        LoadProduct uuid ->
            ( model, fetchProduct uuid )

        PreviousPage ->
            ( {model | pageNumber = model.pageNumber - 1}, fetchProducts (model.pageNumber - 1)  )

        NextPage ->
            ( {model | pageNumber = model.pageNumber + 1}, fetchProducts (model.pageNumber + 1) )

        AddToCart uuid ->
            ( model, addToCart (CartChange Add uuid) )

        GotProducts result ->
            case result of
                Ok pp ->
                    ( { model | content = Just (Products pp), error = "" }, Cmd.none )
                Err e ->
                    ( { model | content = Nothing, error = toString e }, Cmd.none )

        GotProduct result ->
            case result of
                Ok p ->
                    ( { model | content = Just (Product p), error = "" }, Cmd.none )
                Err e ->
                    ( { model | content = Nothing, error = toString e }, Cmd.none )

        AddedToCart result ->
            case result of
                Ok c ->
                    ( { model | cart = c, error = "" }, Cmd.none )
                Err e ->
                    ( { model | cart = [], error = toString e }, Cmd.none )


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
    div [class "mdl-layout mdl-layout--fixed-header"] [
        header [class "mdl-layout__header mdl-layout__header--waterfall custom-header"] [
            div [class "mdl-layout__header-row custom-header-row"] [
                span [class "mdl-layout__title"] [ text "Event Thingy Store" ]
                , div [class "mdl-layout-spacer"] []
--                , div [] [ renderCart model.cart ]
                , div [] [ text model.error ]
                , button [ class "mdl-button mdl-button--raised mdl-button--accent", onClick LoadProducts ] [ text "load products" ]
            ]
        ]
        , div [class "mdl-layout__content", id "main"] [
                div [] [ renderContent model ]
            ]
        ]

renderCart : Cart -> Html msg
renderCart cart =
    text ((String.fromInt (itemsInCart cart)) ++ " items in the cart")

itemsInCart : Cart -> Int
itemsInCart cart =
    foldl (+) 0 (List.map (\e -> 1) cart)

renderContent : Model -> Html Msg
renderContent model =
    case model.content of
        Just (Products pp) -> renderProducts pp model.pageNumber
        Just (Product p) -> renderProductDetail p
        Nothing -> text "" 

renderProducts : Products -> Int -> Html Msg
renderProducts lst pageNumber =
    div [] [
        div [class "mdl-cell"] [
            button [ class "mdl-button mdl-button--raised", onClick PreviousPage ] [ text "previous" ]
            , span [class "mdl-title custom-page-number"] [text (" Page " ++ String.fromInt (pageNumber + 1) ++ " ")]
            , button [ class "mdl-button mdl-button--raised", onClick NextPage ] [ text "next" ]
        ]
        , ul [class "product-list mdl-list"] (List.map (\l -> renderProduct l ) lst )
    ]

renderProduct : ProductOverview  -> Html Msg
renderProduct product =
    li [class "mdl-list__item mdl-list__item--two-line onclick", onClick (LoadProduct product.uuid)]
    [ span [ class "mdl-list__item-primary-content" ] [ span [] [text product.title], span [class "mdl-list__item-sub-title"] [ text "price: ", text (String.fromFloat product.price) ] ]
    , span [ class "mdl-list__item-secondary-content" ] [ button [ class "mdl-button mdl-button--accent mdl-list__item-secondary-action", disabled False ] [ text "add to cart" ] ]
--            , text (String.fromFloat (Maybe.withDefault 0.0 product.price))
            ]


renderProductDetail : ProductDetail  -> Html Msg
renderProductDetail product =
    ol []
    [ li [] [text product.uuid ]
    , li [] [text product.title ]
    , li [] [text product.description ]
    , li [] [text product.longtext ]
    , li [] [text product.category ]
    , li [] [text product.smallImageURL ]
    , li [] [text product.largeImageURL ]
    , li [] [text (String.fromFloat (Maybe.withDefault 0.0 product.price)) ]
    , button [ onClick (AddToCart product.uuid) ] [ text "add to cart" ]
    ]

fetchProducts : Int -> Cmd Msg
fetchProducts pageNumber=
    Http.get
        { url = "http://localhost:8080/products?sort=price&page=" ++ String.fromInt pageNumber
        , expect = Http.expectJson GotProducts productsDecoder
        }

fetchProduct : String -> Cmd Msg
fetchProduct uuid =
    Http.get
        { url = "http://localhost:8080/product?uuid=" ++ uuid
        , expect = Http.expectJson GotProduct productDetailDecoder
        }

addToCart : CartChange -> Cmd Msg
addToCart cartChange =
    Http.riskyRequest
        { method = "POST"
        , headers = []
        , url = "http://localhost:8080/cart"
        , body = Http.jsonBody (encodeCartChange cartChange)
        , expect = Http.expectJson AddedToCart cartDecoder
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


productsDecoder : Decoder Products
productsDecoder = list productDecoder


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
cartItemDecoder = map2 CartItem
  (field "product"  productDetailDecoder)
  (field "quantiy"  int)
