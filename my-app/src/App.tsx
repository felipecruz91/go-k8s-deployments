import React, { useState, useEffect }  from 'react';
import './App.css';
import axios from "axios";

export default function App() {

  const [data, setData] = useState<APIResponse>({Namespace: "x",Deployments:[]});

  useEffect(() => {
    axios
      .get("http://localhost:8090/")
      .then(res => {
        // console.log(result.data); 
        setData(res.data)}
        );
  });

  return (
    <div className="App">
      <header className="App-header">
        <img src="https://freeicons.io/laravel/public/uploads/icons/png/20090363691548218201-128.png" alt="kubernetes-logo"/>
        <p>
           <code>K8s</code> deployments
          
        </p>
      <div> 
        Namespace: <b>{data.Namespace}</b> 
        <ul>
        {data.Deployments.map(item => (
          <li key={item.Name}>
            {item.Name} (Replicas: {item.Replicas})
          </li>
        ))}
      </ul>
      </div>

      </header>
      <div>
    </div>
    </div>
  );
}

interface APIResponse {
  Namespace: string;
  Deployments: Deployment[]
}

interface Deployment{
  Name: string;
  Replicas: Int16Array;
}
