import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css'

// base url for backend
const API_URL = "http://localhost:8080";

function App() {
  const [token, setToken] = useState(null);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [items, setItems] = useState([]);
  
  // view state: 'login' or 'shop'
  const [view, setView] = useState("login");

  useEffect(() => {
    // if we are in shop mode, fetch items
    if (view === "shop") {
      fetchItems();
    }
  }, [view]);

  const fetchItems = async () => {
    try {
      // getting items from backend
      const res = await axios.get(`${API_URL}/items`);
      setItems(res.data);
    } catch (err) {
      console.log("error fetching items");
    }
  };

  const handleLogin = async () => {
    try {
      const res = await axios.post(`${API_URL}/users/login`, {
        username: username,
        password: password
      });
      // saving token
      setToken(res.data.token);
      localStorage.setItem("user_token", res.data.token);
      localStorage.setItem("user_id", res.data.user_id);
      setView("shop");
    } catch (err) {
      window.alert("Invalid username/password");
    }
  };

  const addToCart = async (itemId) => {
    try {
      await axios.post(`${API_URL}/carts`, 
        { item_id: itemId },
        { headers: { token: token } }
      );
      // using toast/alert as requested
      window.alert("Item added to cart!");
    } catch (err) {
      console.error(err);
      window.alert("Failed to add item");
    }
  };

  const showCart = async () => {
    try {
      // logic to get carts
      const res = await axios.get(`${API_URL}/carts`);
      // filter for this user (simple logic for frontend)
      // in real app backend should filter
      const myId = parseInt(localStorage.getItem("user_id"));
      const myCart = res.data.find(c => c.user_id === myId);
      
      if(myCart && myCart.Items) {
         let msg = "Cart Items:\n";
         myCart.Items.forEach(i => msg += `- ${i.name} (ID: ${i.id})\n`);
         window.alert(msg);
      } else {
        window.alert("Cart is empty");
      }
    } catch (err) {
      window.alert("Error fetching cart");
    }
  };

  const showOrders = async () => {
     try {
       const res = await axios.get(`${API_URL}/orders`);
       const myId = parseInt(localStorage.getItem("user_id"));
       const myOrders = res.data.filter(o => o.user_id === myId);
       
       let msg = "Order IDs:\n";
       myOrders.forEach(o => msg += `#${o.id}\n`);
       window.alert(msg);
     } catch(err) {
       window.alert("Error fetching orders");
     }
  };

  const checkout = async () => {
    try {
       // getting cart id first
       const res = await axios.get(`${API_URL}/carts`);
       const myId = parseInt(localStorage.getItem("user_id"));
       const myCart = res.data.find(c => c.user_id === myId);

       if(!myCart) {
         window.alert("No cart found");
         return;
       }

       await axios.post(`${API_URL}/orders`, 
         { cart_id: myCart.id },
         { headers: { token: token } }
       );
       
       setView("shop"); // refresh or stay
       // toast message as requested
       window.alert("Order successful");
    } catch(err) {
      window.alert("Checkout failed");
    }
  };

  // Render Login
  if (view === "login") {
    return (
      <div className="flex items-center justify-center h-screen bg-gray-100">
        <div className="p-8 bg-white shadow-lg rounded-lg">
          <h1 className="text-2xl font-bold mb-4">Login</h1>
          <input 
            className="border p-2 w-full mb-4" 
            placeholder="Username" 
            onChange={e => setUsername(e.target.value)}
          />
          <input 
            className="border p-2 w-full mb-4" 
            type="password"
            placeholder="Password" 
            onChange={e => setPassword(e.target.value)}
          />
          <button 
            className="bg-blue-500 text-white p-2 w-full rounded"
            onClick={handleLogin}
          >
            Login
          </button>
        </div>
      </div>
    );
  }

  // Render Shop
  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Shop Items</h1>
        <div className="space-x-4">
           <button onClick={checkout} className="bg-green-500 text-white px-4 py-2 rounded">Checkout</button>
           <button onClick={showCart} className="bg-yellow-500 text-white px-4 py-2 rounded">Cart</button>
           <button onClick={showOrders} className="bg-gray-500 text-white px-4 py-2 rounded">Order History</button>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4">
        {items.map(item => (
          <div key={item.id} className="border p-4 rounded shadow hover:shadow-lg transition">
            <h2 className="text-xl font-semibold">{item.name}</h2>
            <p className="text-gray-500">{item.status}</p>
            <button 
              onClick={() => addToCart(item.id)}
              className="mt-4 bg-blue-500 text-white px-4 py-1 rounded text-sm"
            >
              Add to Cart
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

export default App;