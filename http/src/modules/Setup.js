import React from "react";
import { Link, browserHistory } from "react-router";

import { POST } from "../utils/ajax";

// Default placeholder DB strings.
const DB_STRINGS = {
	"sqlite3": "database.sqlite",
	"mysql": "user:password@/dbname?charset=utf8&parseTime=True&loc=Local",
	"postgres": "user=soandso password=secret dbname=scarecrow",
}

export default React.createClass({
	getInitialState() {
		return {
			dbType: "sqlite3",
			dbString: "database.sqlite",
			adminName: "",
			adminPassword1: "",
			adminPassword2: "",
		}
	},

	handleSubmit(e) {
		e.preventDefault();

		// Validate things.
		if (this.state.dbString.length === 0) {
			window.alert("The database connection string is required.");
		}
		else if (this.state.adminName.length === 0) {
			window.alert("The admin username is required.");
		}
		else if (this.state.adminPassword1.length === 0) {
			window.alert("You must enter a password for the admin user.");
		}
		else if (this.state.adminPassword1 !== this.state.adminPassword2) {
			window.alert("The admin passwords do not match.");
		}
		else {
			POST("/v1/admin/setup", this.state, function(data) {
				serverSettings.initialized = true;
				serverSettings.loggedIn = true;
				browserHistory.push("/");
			}, function(errMsg) {
				window.alert(errMsg);
			});
		}
	},

	handleDbTypeChange(e) {
		// Handle them changing the database driver selection.
		// If they hadn't modified the connection string away from its default
		// value, then set it to be the new default for the new driver.
		var stillDefaultString = this.state.dbString === DB_STRINGS[this.state.dbType];
		var newType = e.target.value;
		this.setState({dbType: newType});

		if (stillDefaultString) {
			this.setState({dbString: DB_STRINGS[newType]});
		}
	},

	handleDbStringChange(e) {
		this.setState({dbString: e.target.value});
	},

	handleAdminNameChange(e) {
		this.setState({adminName: e.target.value});
	},

	handleAdminPassword1Change(e) {
		this.setState({adminPassword1: e.target.value});
	},

	handleAdminPassword2Change(e) {
		this.setState({adminPassword2: e.target.value});
	},

	render() {
		if (serverSettings.initialized === true) {
			return this.alreadyInitialized();
		}
		return (
			<div>
				<h1 className="page-header">Setup Scarecrow</h1>

				<div className="col-sm-8">
					<p>
						Welcome to Scarecrow! In just a few minutes you'll be able
						to start creating your own RiveScript bots!
					</p>

					<form onSubmit={this.handleSubmit}>
						<h2>Database Setup</h2>

						<p>
							Enter your preferred database settings. The simplest option
							is to use SQLite3, which creates a database as a local file
							on disk.
						</p>

						<p>
							Supported databases are: MySQL, PostgreSQL, and SQLite3.
						</p>

						<div className="form-group">
							<label for="db-type">Database Type</label>
							<select
								name="db-type"
								className="form-control"
								selected={this.state.dbString}
								onChange={this.handleDbTypeChange}>
									<option value="sqlite3">SQLite3</option>
									<option value="postgres">PostgreSQL</option>
									<option value="mysql">MySQL</option>
							</select>
						</div>

						<div className="form-group">
							<label for="db-string">Connection String</label>
							<input type="text"
								className="form-control"
								value={this.state.dbString}
								onChange={this.handleDbStringChange} />
						</div>

						<h2>Admin Account</h2>

						<p>
							Set up an administrator account for logging into this
							web app and configuring your bots.
						</p>

						<div className="form-group">
							<label for="admin-name">Admin Username</label>
							<input type="text"
								name="admin-name"
								className="form-control"
								placeholder="admin"
								value={this.state.adminName}
								onChange={this.handleAdminNameChange} />
						</div>

						<div className="form-group">
							<label for="admin-password1">Admin Password</label>
							<input type="password"
								name="admin-password1"
								className="form-control"
								placeholder="correct horse battery staple"
								value={this.state.adminPassword1}
								onChange={this.handleAdminPassword1Change} />
							<input type="password"
								name="admin-password2"
								className="form-control"
								placeholder="(confirm)"
								value={this.state.adminPassword2}
								onChange={this.handleAdminPassword2Change} />
						</div>

						<div className="form-group">
							<button type="submit" className="btn btn-primary">Continue</button>
						</div>
					</form>
				</div>
			</div>
		);
	},

	// The view sent when the app has already been set up.
	alreadyInitialized() {
		return (
			<div>
				<h1 className="page-header">Setup Scarecrow</h1>

				<p>
					This app has already been configured. Please go to{' '}
					<Link to="/login">log in</Link> instead.
				</p>
			</div>
		)
	}
});
