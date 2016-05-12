import React from "react";
import { Link, browserHistory } from "react-router";

import { POST } from "../utils/ajax"

export default React.createClass({
	getInitialState() {
		return {
			username: "",
			password: "",
		}
	},

	handleSubmit(e) {
		e.preventDefault();

		// Validate things.
		if (this.state.username.length === 0) {
			window.alert("The admin username is required.");
		}
		else if (this.state.password.length === 0) {
			window.alert("The password is required.");
		}
		else {
			POST("/v1/admin/auth", this.state, function(data) {
				if (data.status === "ok") {
					serverSettings.loggedIn = true;
					browserHistory.push("/");
				}
				else {
					window.alert(data.message);
				}
			}, function(errMsg) {
				window.alert(errMsg);
			});
		}
	},

	handleUsernameChange(e) {
		this.setState({username: e.target.value});
	},

	handlePasswordChange(e) {
		this.setState({password: e.target.value});
	},

	render() {
		return (
			<div>
				<h1 className="page-header">Please log in.</h1>

				<div className="col-sm-8">
					<p>
						Please log in with your admin username and password.
					</p>

					<form onSubmit={this.handleSubmit}>

						<div className="form-group">
							<label for="admin-name">Admin Username</label>
							<input type="text"
								name="username"
								className="form-control"
								value={this.state.username}
								onChange={this.handleUsernameChange} />
						</div>

						<div className="form-group">
							<label for="admin-password1">Admin Password</label>
							<input type="password"
								name="password"
								className="form-control"
								value={this.state.password}
								onChange={this.handlePasswordChange} />
						</div>

						<div className="form-group">
							<button type="submit" className="btn btn-primary">Continue</button>
						</div>
					</form>
				</div>
			</div>
		);
	}
});
