const url = "http://localhost:3000/items";

function renderList() {
    fetch(url).then(
        function(response) {
            return response.json();;
        }
    ).then(
        function(resp) {
            const time =  new Date().toLocaleTimeString();
        console.log(resp, time);
        console.log("rendering list: actually");
        const mappedItems = resp.data.map(item => {
            return `
            <li>
               <h3> ${item.name} </h3>
                <p> ${item.description} </p>
                <p> ${item.quantity} </p>
                <button onclick='showMenu("${item.id}", "${item.name}", "${item.description}", "${item.quantity}")'> Expand </button>
                <div id="${item.id}" hidden="true">
                    <input type="text" name="comment" placeholder="comment">
                    <button onclick='deleteItem("${item.id}")'> Delete </button>
                    <br>
                    <input type="text" name="edit_name" placeholder="name">
                    <input type="text" name="edit_description" placeholder="description">
                    <input type="text" name="edit_quantity" placeholder="quantity">
                    <button onclick='updateItem("${item.id}")'> Update </button>
                </div>
            </li>
            `;
        })

        document.getElementById("inventory").innerHTML = mappedItems.join("");
    }).catch(function(error) {
        console.log(error);
    });

}

function createItem() {
    const quantity = document.querySelector("input[name='quantity']").value;
    const name = document.querySelector("input[name='name']").value;
    const description = document.querySelector("input[name='description']").value;
    console.log(quantity, name, description);
    let intQuantity = parseInt(quantity);
    let payload = {
        name: name,
        description: description,
        quantity: intQuantity
    }

    fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    }).then(
        function(response) {
            return response.json();
        }
    ).then(
        function(resp) {
            console.log(resp);
            location.reload();
        } 
    ).catch(function(error) {
        console.log("debug error",error);
    })

}

function deleteItem(id) {
    const comment = document.querySelector("input[name='comment']").value;
    const payload = {
        "comment": comment
    }
    console.log("in here", id)
    fetch(url + "/" + id, {
        method: "DELETE",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    }).then(
        function(response) {
            return response.json();
        }
    ).then(
        function(resp) {
            console.log(resp);
            location.reload();
        }
    ).catch(function(error) {
        console.log(error);
    })
}

function showMenu(id, name, description, quantity) {
    // const id = obj.id ? obj.id : obj;
    console.log("in here show menu", id);
    // document.getElementById(id).style.display = "block";
    document.getElementById(id).hidden = !document.getElementById(id).hidden;
    document.querySelector("input[name='edit_name']").value = name;
    document.querySelector("input[name='edit_description']").value = description;
    document.querySelector("input[name='edit_quantity']").value = quantity;
    
}

function restoreItem(id) {
    console.log("in here", id)
    fetch(url + "/" + id, {
        method: "PUT",
    }).then(
        function(response) {
            return response.json();
        }
    ).then(
        function(resp) {
            console.log(resp);
            location.reload();
        }
    ).catch(function(error) {
        console.log(error);
    })
}

function showArchive() {
    console.log("in here show archive");
    fetch(url +"?"+ new URLSearchParams({
        "status": true
    })).then(function(response) {
        return response.json();
    } ).then(function(resp) {
        console.log(resp);
        const mappedItems = resp.data.map(item => {
            return `
            <li>
                <h3> ${item.name} </h3>
                <p> ${item.description} </p>
                <p> Comment: ${item.comment} </p>
                <p> ${item.quantity} </p>

                <button id="${item.id}" onclick='restoreItem("${item.id}")'> Restore </button>
            </li>
            `;
        })
        document.getElementById("archive").innerHTML = mappedItems.join("");

    }).catch(function(error) {
        console.log(error);
    }
    )
}

function updateItem(id) {
    const quantity = document.querySelector("input[name='edit_quantity']").value;
    const name = document.querySelector("input[name='edit_name']").value;
    const description = document.querySelector("input[name='edit_description']").value;
    console.log(quantity, name, description);
    let intQuantity = parseInt(quantity);
    let payload = {
        name: name,
        description: description,
        quantity: intQuantity
    }
    fetch(url + "/" + id, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    }).then(
        function(response) {
            return response.json();
        }
    ).then(
        function(resp) {
            console.log(resp);
            location.reload();
        }
    ).catch(function(error) {
        console.log(error);
    })
}

// const variable = document.querySelector(".submit");
// variable.addEventListener("click", createItem);

// i can hear you

renderList();

