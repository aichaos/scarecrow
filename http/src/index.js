import React from "react"
import { render } from "react-dom"
import { Router, Route, browserHistory, IndexRoute } from "react-router"

// Our pages.
import App from "./modules/App"
import Home from "./modules/Home"
import Setup from "./modules/Setup"

render((
	<Router history={browserHistory}>
		<Route path="/" component={App}>
			<IndexRoute component={Home}/>

			<Route path="/setup" component={Setup}/>
		</Route>
	</Router>
), document.getElementById("app"));
