import React from "react";

import TopNav from "./TopNav"
import LeftNav from "./LeftNav"

export default React.createClass({
	render() {
		return (
			<div>
				<TopNav/>

				<div className="container-fluid">
					<div className="row">
						<LeftNav/>
						<div className="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
							{this.props.children}
						</div>
					</div>
				</div>
			</div>
		);
	}
});
