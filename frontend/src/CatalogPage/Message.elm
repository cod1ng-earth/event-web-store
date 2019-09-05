module CatalogPage.Message exposing (SubMsg(..))

import Http
import Catalog


type SubMsg
    = GotProducts (Result Http.Error Catalog.CatalogResponse)
    | LoadProducts
    | PreviousPage
    | NextPage
    | SortByUuid
    | SortByPrice
    | SortByName
    | DisableFilterByPrefix
    | EnableFilterByPrefix
    | GoToPage String
    | SetFilterPrefix String
