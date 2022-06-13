# shopify

github.com/kasumusof/shopify

## How to use the app

### front end
###### Hitting the root (/) would serve the a html page.
###### A script would query all the items in the inventory and render them on the page. 
###### A form exists at the top of the page to add an item to the inventory. On a successful add, the item is rendered on the page.
###### Each item has a button "Expand" under them that would show 2 forms when clicked.A form to delete (comment is reqiured) the item and a form to update the item.

###### A button exists on the bottom of the page "Show Archive". This button query all the archived items in the inventory and their comments (on deletion).


### backend

#### /items
```post```
###### Allows a user to create an item in the inventory
##### Body params
###### ```name``` string required
###### ```description``` string required
###### ```quantity``` integer required
##### Query params
###### ```nil```


##### Responses
###### ```201``` status 201 response
###### ```400``` status 400 response

##### Response Body Example
```
{
  message : string
  data : {
    id : uuid
    name : string
    description : string
    created_at : timestamp
    updated_at : timestamp
  }
}
```


#### /items
```get```
###### Retrieve all items in the inventory

##### Body params
###### ```nil```

##### Query params
###### ```status``` string optional
###### When this param is passed it returns the archived (deleted) items. When Empty, it returns the items in the inventory

##### Responses
###### ```201``` status 201 response
###### ```500``` status 500 response

##### Response Body Example
```
{
  message : string
  data : [
    {
      id : uuid
      name : string
      description : string
      created_at : timestamp
      updated_at : timestamp
      comment : string
      deleted_at : timestamp
    }
  ]
}
```


#### /items/:id
```delete```
###### Delete (archive) an item in the inventory

##### Body params
###### ```comment``` string required

##### Path params
###### ```id``` uuid required

##### Query params
###### ```nil```

##### Responses
###### ```201``` status 201 response
###### ```400``` status 400 response
###### ```500``` status 500 response

##### Response Body Example
```
{
  message : string
  data : uuid
}
```


#### /items/:id
```patch```
###### Update an item in the inventory

##### Body params
###### ```nil```

##### Path params
###### ```id``` uuid required

##### Query params
###### ```nil```

##### Responses
###### ```201``` status 201 response
###### ```400``` status 400 response
###### ```500``` status 500 response

##### Response Body Example
```
{
  message : string
  data : uuid
}
```

#### /items/:id
```put```
###### Unarchive an item.

##### Body params
###### ```nil```

##### Path params
###### ```id``` uuid required

##### Query params
###### ```nil```

##### Responses
###### ```201``` status 201 response
###### ```400``` status 400 response
###### ```500``` status 500 response

##### Response Body Example
```
{
  message : string
  data : uuid
}
```

