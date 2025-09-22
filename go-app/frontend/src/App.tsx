// import { useState } from "react";
// import logo from "./assets/images/logo-universal.png";
// import "./App.css";
// import { Greet } from "../wailsjs/go/main/App";

function App() {
  // const [resultText, setResultText] = useState(
  //   "Please enter your name below 👇",
  // );
  // const [name, setName] = useState("");
  // const updateName = (e: any) => setName(e.target.value);
  // const updateResultText = (result: string) => setResultText(result);

  // function greet() {
  //   Greet(name).then(updateResultText);
  // }

  return (
    <div className="p-6">
      <select className="select w-full">
        <option disabled selected>
          Pick a color
        </option>
        <option>Crimson</option>
        <option>Amber</option>
        <option>Velvet</option>
      </select>
    </div>
  );
}

export default App;
