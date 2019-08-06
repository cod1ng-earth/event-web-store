module Main exposing (main)

-- https://package.elm-lang.org/packages/NoRedInk/elm-sweet-poll/latest/
-- https://package.elm-lang.org/packages/simonh1000/elm-jwt/latest/
-- https://github.com/simonh1000/elm-jwt/blob/master/examples/node/src/Main.elm

import Browser
import Browser.Events exposing (onVisibilityChange, Visibility(..))
import Html exposing (Html, button, div, text, h1, ol, ul, li, a)
import Html.Events exposing (onClick)
import Html.Attributes exposing (href)
import Time
import Task
import Random
import Http
import Delay
import Json.Decode exposing (Decoder, map3, map8, float, string, field, list, decodeString)


type alias Model =
    { error : String
    , content : Maybe Content
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

init : () -> ( Model, Cmd Msg )
init _ =
    ( { error = ""
      , content = Nothing
      }
    , Cmd.batch [ ]
    )


type Msg
    = LoadProducts
    | LoadProduct String
    | GotProducts (Result Http.Error Products)
    | GotProduct (Result Http.Error ProductDetail)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of

        LoadProducts ->
            ( model, fetchProducts )

        LoadProduct uuid ->
            ( model, fetchProduct uuid )

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
        [ h1 [] [ text "Hello, world." ]
        , div []
            [ button [ onClick LoadProducts ] [ text "load products" ]
            , button [ onClick (LoadProduct "0002630c-718d-4ffb-8989-a760cdf69c26") ] [ text "load 0002630c-718d-4ffb-8989-a760cdf69c26" ]
            , div [] [ text model.error ]
            , div [] [ renderContent model.content ]
            ]
        ]

renderContent : Maybe Content -> Html Msg
renderContent content =
    case content of
        Just (Products pp) -> renderProducts pp
        Just (Product p) -> renderProductDetail p
        Nothing -> text "" 

renderProducts : Products -> Html Msg
renderProducts lst =
    ul []
        (List.map (\l -> li [] [ renderProduct l ]) lst )

renderProduct : ProductOverview  -> Html Msg
renderProduct product =
    div []
            [ button [ onClick (LoadProduct product.uuid) ] [ text "more details!" ]
            , text product.title
            ]


renderProductDetail : ProductDetail  -> Html msg
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
    ]

fetchProducts : Cmd Msg
fetchProducts =
    Http.get
        { url = "http://localhost:8080/all-products"
        , expect = Http.expectJson GotProducts productsDecoder
        }

fetchProduct : String -> Cmd Msg
fetchProduct uuid =
    Http.get
        { url = "http://localhost:8080/product?uuid=" ++ uuid
        , expect = Http.expectJson GotProduct productDetailDecoder
        }


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