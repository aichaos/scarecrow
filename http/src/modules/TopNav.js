import React from "react";
import { Link, browserHistory } from "react-router";

import { POST } from "../utils/ajax"

export default React.createClass({
	handleLogout(e) {
		e.preventDefault();

		POST("/v1/admin/deauth", "{}", function(data) {
			serverSettings.loggedIn = false;
			browserHistory.push("/login");
		}, function(errMsg) {
			window.alert(errMsg);
		});
	},

	render() {
		if (serverSettings.loggedIn === true) {
			return (
				<ul className="nav navbar-nav navbar-right">
					<li><Link to="/">Dashboard</Link></li>
					<li><Link to="/setup">Setup</Link></li>
					<li><Link to="/settings">Settings</Link></li>
					<li><a href="#" onClick={this.handleLogout}>Log Out</a></li>
					<li><Link to="/help">Help</Link></li>
					<li></li>
				</ul>
			);
		}
		else if (serverSettings.initialized === false) {
			return (
				<ul className="nav navbar-nav navbar-right">
					<li><Link to="/setup">Setup</Link></li>
					<li><Link to="/help">Help</Link></li>
				</ul>
			);
		}
		else {
			return (
				<ul className="nav navbar-nav navbar-right">
					<li><Link to="/login">Log In</Link></li>
					<li><Link to="/help">Help</Link></li>
				</ul>
			);
		}
	}
});
