// import config from '../../config'
import Content from './content'
import WSController from './wscontroller'
const port = '8100'

export default class App extends React.PureComponent {
	render() {
		return (
			<div>
				<WSController url={`ws://localhost:${port}/`}>
					<Content />
				</WSController>
			</div>
		)
	}
}
