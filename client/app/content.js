import Board from './board'
import Form from './form'

export default class Content extends React.PureComponent {
	render() {
		return (
			<div>
				<Form {...this.props} />
				<Board {...this.props} />
			</div>
		)
	}
}