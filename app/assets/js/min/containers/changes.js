!function e(t,a,n){function r(l,i){if(!a[l]){if(!t[l]){var o="function"==typeof require&&require;if(!i&&o)return o(l,!0);if(c)return c(l,!0);var u=new Error("Cannot find module '"+l+"'");throw u.code="MODULE_NOT_FOUND",u}var s=a[l]={exports:{}};t[l][0].call(s.exports,function(e){var a=t[l][1][e];return r(a?a:e)},s,s.exports,e,t,a,n)}return a[l].exports}for(var c="function"==typeof require&&require,l=0;l<n.length;l++)r(n[l]);return r}({1:[function(e,t,a){$(document).ready(function(){$("#menu-containers").addClass("active")});var n=React.createClass({displayName:"TableRow",render:function(){var e=this.props.content,t="";switch(e.Kind){case 0:t="Modify";break;case 1:t="Add";break;case 2:t="Delete"}return React.createElement("tr",{key:this.props.index},React.createElement("td",{className:"data-name"},t),React.createElement("td",{className:"data-name"},e.Path))}}),r=React.createClass({displayName:"Table",getInitialState:function(){return{data:[]}},componentDidMount:function(){var e=this,t=$("#container-id").val(),a=$("#client-id").val();a=a?"?client="+a:"",app.func.ajax({type:"GET",url:"/api/container/changes/"+t+a,success:function(t){e.setState({data:t})}})},render:function(){var e;return e=this.state.data.length>0?this.state.data.map(function(e,t){return React.createElement(n,{index:t,content:e})}):React.createElement(n,{index:0,content:{Kind:-1,Path:"There is no changed file."}}),React.createElement("table",{className:"table table-striped table-hover"},React.createElement("thead",null,React.createElement("tr",null,React.createElement("th",null),React.createElement("th",null,"Path"))),React.createElement("tbody",null,e))}});React.render(React.createElement(r,null),document.getElementById("data"))},{}]},{},[1]);