import logo from './logo.svg';
import './App.css';
import React from "react";

const baseURL = 'http://localhost:3001/contract'
export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      contract: '',
      abi: '',
      methods: null,
      init: true,
      error: '',
      selectedMethod: '',
      callData: {},
      blockNumber: '',
      result: [],
      rememberABI: false
    };
  }

  handleChangeContract(e) {
    this.setState({ contract: e.target.value });
  }

  handleChangeABI(e) {
    this.setState({ abi: e.target.value });
  }

  initLayout() {
    return (
      <div className="container container-init-state">
        <div className="contract contract-address">
          <div className="label">Contract Address</div>
          <input onFocus={(e) => {this.setError('')}} onChange={(e) => {this.handleChangeContract(e)}} value={this.state.contract} type="text"/>
        </div>
        <div className="contract contract-abi">
        <div className="label">Contract ABI</div>
          <textarea placeholder="This is optional. If you leave it empty, we will try to look up the ABI in our database and Etherscan..." onChange={(e) => {this.handleChangeABI(e)}} value={this.state.abi}/>
        </div>
        {this.state.error == '' ? '' : <div className="contract-error">Error: {this.state.error}</div>}
        <div className={this.state.rememberABI ? "contract contract-check-box contract-check-box__checked" : "contract contract-check-box"} onClick={(e) => {this.setState({rememberABI: !this.state.rememberABI})}}>
          <div className="check-box"></div>
          <div className="check-box-content">Remember the ABI</div>
        </div>
        <div className="contract contract-submit">
          <button onClick={(e) => {this.verifyAndAccessContract()}}>Submit</button>
        </div>
      </div>
    )
  }

  verifyAndAccessContract() {
    var data = {
      contract: this.state.contract,
      abi: String.raw`${this.state.abi}`,
      rememberABI: this.state.rememberABI
    }
    var url = `${baseURL}/methods`
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data)
    }).then(response => response.json()).then(data => {
      if (data.err) {
        this.setError(data.err)
        return
      }
      if (Array.isArray(data.data) && data.data.length > 0) {
        this.setState({methods: data.data, selectedMethod: data.data[0].name, init: false})
      } else {
        this.setError("cannot get data from server")
        return
      }
    }).catch(err => this.setState({error: err}))
  }

  setError = (err) => {
    this.setState({error: err})
  }


  // contract call
  handleSelectChange = e => {
    this.setState({selectedMethod: e.target.value, callData: {}})
  }

  handleInputChange = (e, name) => {
    var callData = {...this.state.callData}
    callData[name] = e.target.value
    this.setState({callData: callData})
  }

  generateArgs = (args) => {
    return args.map((arg, index) => {
      const name = arg.name
      return (
        <div className="param" key={index}>
          <div className="param-label">
            <div className="name">
              {name == "" ? "input" : name}
            </div>
          </div>
          <div className="param-input">
            <input 
              key={index} 
              name={name} 
              value={this.state.callData[name] || ''}
              onChange={(e) => this.handleInputChange(e, name)}
              placeholder={arg.type}
              onFocus={(e) => {this.setError('')}} 
            />
          </div>
        </div>
      )
    })
  }

  submitContractData() {
    var data = {
      contract: this.state.contract,
      abi: String.raw`${this.state.abi}`,
      method: this.state.selectedMethod,
      params: this.state.callData,
      blockNumber: this.state.blockNumber
    }
    console.log(JSON.stringify(data))
    var url = `${baseURL}/call`
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data)
    }).then(response => response.json()).then(data => {
      if (data.err) {
        this.setError(data.err)
        return
      }
      console.log(data.data)
      this.setState({result: data.data})
    }).catch(err => this.setError(err))
  }

  handleInputBlockNumberChange = (e) => {
    this.setState({blockNumber: e.target.value})
  }

  callLayout() {
    const options = this.state.methods.map((method, index) => <option key={index} value={method.name}>{method.name}</option>);
    const params = this.state.methods.map((method, index) => {
      if (method.name == this.state.selectedMethod) {
        return (
          <div key={index}>
            {Array.isArray(method.arguments) ? this.generateArgs(method.arguments) : ''}
          </div>
        )
      }
    });
    return (
      <div className="container container-call">
        <div className="contract contract-method">
          <div className="label">Methods</div>
          <select 
            value={this.state.selectedMethod} 
            onChange={this.handleSelectChange} 
            name="methods" 
            id="methods"
          >
            {options}
          </select>
        </div>
        <div className="contract contract-params">
          {params}
          <div className="param">
            <div className="param-label">
              <div className="name">
                block nummber
              </div>
            </div>
            <div className="param-input">
              <input 
                name={'blockNumber'} 
                value={this.state.blockNumber}
                onChange={(e) => this.handleInputBlockNumberChange(e)}
                placeholder={'default is latest block'}
                onFocus={(e) => {this.setError('')}} 
              />
            </div>
          </div>
        </div>
        {this.state.error == '' ? '' : <div className="contract-error">Error: {this.state.error}</div>}
        <div className="result-label">Result</div>
        <div className="contract contract-result">
          {this.state.result.map((data, index) => {
            return <div className={"contract-result_element"} key={index}> [{index+1}]: {JSON.stringify(data)}</div>
          })}
        </div>
        <div className="contract contract-submit">
          <button onClick={(e) => {this.submitContractData()}}>Submit</button>
        </div>
      </div>
    )
  }

  reset() {
    this.setState({
      contract: '',
      abi: '',
      methods: null,
      init: true,
      error: '',
      selectedMethod: '',
      callData: {},
      blockNumber: '',
      result: []
    })
  }
 
  render() {
    return (
      <div className="wrapper">
        <div className="title" onClick={(e) => this.reset()}>Contract Caller</div>
        <div className="content">
          {this.state.init ? this.initLayout() : this.callLayout()}
        </div>
      </div>
    )
  }
}
