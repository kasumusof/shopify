let url = "http://localhost:8080/items";

fetch(url).then(
    function(response) {
        return response.json();;
    }
).then(
    function(resp) {const inventory = document.getElementById("inventory");

    console.log(resp)
    const data = resp.message ? resp.data : []
    console.log(data);
    
    const mappedItems = data.map(item => {
        return `<li>
           <h3> ${item.name} </h3>
            <p> ${item.description} </p>
            <p> ${item.quantity} </p>
        </li>`;
    })
    
    inventory.innerHTML += mappedItems.join("");}
).catch(function(error) {
    console.log(error);
});

