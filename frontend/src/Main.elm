module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1, ol, ul, li, a)
import Html.Events exposing (onClick)
import Html.Attributes exposing (href, class)
import List exposing (foldl)
import Time
import Task
import Random
import Http
import Delay
import Json.Decode exposing (Decoder, map2, map3, map8, int, float, string, field, list, decodeString)
import Json.Encode


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
  , price:         Float
  }

type alias Products = List ProductOverview
type alias ProductOverview =
  { uuid             : String
  , title            : String
  , price            : Float
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
    div []
        [ h1 [ class "dings" ] [ text "Hello, world." ]
        , div []
            [ button [ onClick LoadProducts ] [ text "show products" ]
            , div [] [ renderCart model.cart ]
            , div [] [ text model.error ]
            , div [] [ renderContent model ]
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
    div []
    [ button [ onClick PreviousPage ] [ text "previous" ]
    , text (" Page " ++ String.fromInt (pageNumber + 1) ++ " ")
    , button [ onClick NextPage ] [ text "next" ]
    , ul [] (List.map (\l -> li [] [ renderProduct l ]) lst )
    ]

renderProduct : ProductOverview  -> Html Msg
renderProduct product =
    div []
            [ button [ onClick (LoadProduct product.uuid) ] [ text "more details!" ]
            , button [ onClick (AddToCart product.uuid) ] [ text "add to cart" ]
            , text " "
            , text product.title
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
    , li [] [text (String.fromFloat product.price) ]
    , button [ onClick (AddToCart product.uuid) ] [ text "add to cart" ]
    ]

fetchProducts : Int -> Cmd Msg
fetchProducts pageNumber=
    Http.get
        { url = "http://localhost:8080/products?page=" ++ String.fromInt pageNumber
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
    Http.post
        { url = "http://localhost:8080/cart"
        , body = Http.jsonBody (encodeCartChange cartChange)
        , expect = Http.expectJson AddedToCart cartDecoder
        }

encodeCartChange : CartChange -> Json.Encode.Value
encodeCartChange cc  =
    Json.Encode.object
        [ ( "action", Json.Encode.string (cartActionToString cc.action) )
        , ( "uuid", Json.Encode.string cc.uuid )
        ]

cartActionToString : CartAction -> String
cartActionToString c =
    case c of
      Add -> "add"
      Remove -> "add"

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
productDecoder = map3 ProductOverview
  (field "uuid"          string)
  (field "title"         string)
  (field "price"         float)


productDetailDecoder : Decoder ProductDetail
productDetailDecoder = map8 ProductDetail
  (field "uuid"          string)
  (field "title"         string)
  (field "description"   string)
  (field "longtext"      string)
  (field "category"      string)
  (field "smallImageURL" string)
  (field "largeImageURL" string)
  (field "price"         float)


cartDecoder : Decoder Cart
cartDecoder = list cartItemDecoder

cartItemDecoder : Decoder CartItem
cartItemDecoder = map2 CartItem
  (field "product"  productDetailDecoder)
  (field "quantiy"  int)
