var App = React.createClass({
    componentWillMount: function () {
    },
    render: function () {

        if (this.loggedIn) {
            return (<LoggedIn />);
        } else {
            return (<Home />);
        }
    }
});