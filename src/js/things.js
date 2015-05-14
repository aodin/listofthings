var ListOfThings = function() {
  'use strict';

  var colors = [
    '#f0585e', // Red
    '#5899d2', // Blue
    '#78c269', // Green
    '#f9a65a', // Orange
    '#9d65aa', // Purple
    '#4fc99d', // Teal
    '#cc6f57', // Magenta
    '#d67eb2', // Lavender
    '#5c68aa', // Blue 2
    '#5ce160', // Toxic Green
    '#d6d67e', // Ugly yellow
  ];

  var SESSION_ID = 'listuserid';
  var WEBSOCKET_URI = 'ws://' + document.URL.split('/', 3)[2] + '/feeds/v1/things';

  var App = Backbone.View.extend({
    el: '#main',
    initialize: function(options) {
      // Save the users

      // Create a new websocket
      this.ws = new WebSocket(WEBSOCKET_URI);
      this.ws.onopen = this.join.bind(this);
      this.ws.onmessage = this.onMessage.bind(this);
      this.ws.onerror = this.onError.bind(this);
      this.ws.onclose = this.leave.bind(this);
    },
    getCookies: function() {
      var c = document.cookie, v = 0, cookies = {};
      if (document.cookie.match(/^\s*\$Version=(?:"1"|1);\s*(.*)/)) {
        c = RegExp.$1;
        v = 1;
      }
      if (v === 0) {
        c.split(/[,;]/).map(function(cookie) {
          var parts = cookie.split(/=/, 2),
            name = decodeURIComponent(parts[0].trimLeft()),
            value = parts.length > 1 ? decodeURIComponent(parts[1].trimRight()) : null;
          cookies[name] = value;
        });
      } else {
        c.match(/(?:^|\s+)([!#$%&'*+\-.0-9A-Z^`a-z|~]+)=([!#$%&'*+\-.0-9A-Z^`a-z|~]*|"(?:[\x20-\x7E\x80\xFF]|\\[\x00-\x7F])*")(?=\s*[,;]|$)/g).map(function($0, $1) {
          var name = $0,
              value = $1.charAt(0) === '"' ? $1.substr(1, -1).replace(/\\(.)/g, "$1") : $1;
          cookies[name] = value;
        });
      }
      return cookies;
    },
    getCookie: function(name) {
      return this.getCookies()[name];
    },
    join: function() {
      this.ws.send(JSON.stringify({session: this.getCookie(SESSION_ID)}));
    },
    onMessage: function(msg) {
      console.log('new message:', msg);
    },
    onError: function(err) {},
    leave: function() {},
    sync: function(method, model) {
      // Check ready state
      console.log('sync:', method, model, model.toJSON(), this.ws.readyState);
      if (this.ws.readyState !== 1) {
        // Return after displaying an error
        $('#errors').prepend(new Error({message: 'Could not connect to server'}).el);
        return;
      }

      var msg = {
        method: method,
        content: model.toJSON()
      };
      this.ws.send(JSON.stringify(msg));
    }
  });

  var module = {};
  module.onready = function() {
    // App needs to know both users and the collection because all socket
    // messages go through it
    var app = new App({collection: new Things(), users: new Users()});

    // Bind to the app's sync
    Backbone.sync = app.sync.bind(app);
  };

  // Self-deleting error messages
  var Error = Backbone.View.extend({
    tagName: 'li',
    initialize: function(options) {
      this.message = options.message;
      var timeout = options.timeout || 4000;

      // Destory self after a timeout
      var self = this;
      setTimeout(function() {self.remove();}, timeout);

      // Render automatically
      return this.render();
    },
    render: function() {
      this.$el.html(this.message);
      return this;
    }
  });

  var List = Backbone.View.extend({
    el: '#things',
    events: {
      'keyup #create-name': 'proxyEnter',
      'click #create': 'createItem',
    },
    initialize: function(options) {
      this.users = options.users;
      _.bindAll(this, 'onmessage', 'onopen', 'onerror', 'onclose');
      this.listenTo(this.collection, 'reset', this.render);
      this.listenTo(this.collection, 'add', this.renderItem);
    },
    onmessage: function(result) {
      var data = JSON.parse(result.data);
      if (!data) {return;}

      // TODO This should be behavior of Backbone.sync
      switch (data.body) {
        case 'init':
          this.collection.reset(data.content);
          break;
        case 'create':
          this.collection.add(data.content);
          break;
        case 'delete':
          var thing = this.collection.get(data.content.id);
          this.collection.remove(thing);
          break;
        case 'update':
          this.collection.add(data.content, {merge: true});
          break;
        case 'users':
        case 'join':
          this.users.add(data.content);
          break;
        case 'left':
          var user = this.users.get(data.content.id);
          this.users.remove(user);
          break;
        case 'error':
          // Error message, display for a timeout
          this.$('#errors').prepend(new Error({message: data.content}).el);
          break;
        default:
          // TODO Error
      }
    },
    // websocket methods have default parameter: data
    onopen: function() {},
    onerror: function() {},
    onclose: function() {},
    proxyEnter: function(e) {
      if (e.keyCode === 13) {this.createItem();}
    },
    createItem: function() {
      var $input = this.$('#create-name');

      // Trim whitespace
      // TODO thing validation
      var name = $.trim($input.val());
      if (!name) {
        this.$('#errors').prepend(new Error({message: 'Empty items cannot be created'}).el);
        return;
      }

      // TODO Create through the collection methods
      var thing = new Thing({name: $input.val()});
      thing.save();

      // Clear the input
      $input.val('');
    },
    render: function() {
      var $list = this.$('ol');
      $list.empty();
      _.each(this.collection.models, function(m) {
        var item = new Item({model: m});
        $list.append(item.render().el);
      }, this);
      return;
    },
    renderItem: function(m) {
      var $list = this.$('ol');
      var item = new Item({model: m});
      $list.append(item.render().el);
    }
  });

  var Item = Backbone.View.extend({
    tagName: 'li',
    template: _.template('<h3><%= name %> <span class="edit"><small>edit</small></span><span class="delete"><small>delete</small></span></h3>'),
    editTemplate: _.template('<div class="input-group"><input type="text" class="form-control" value="<%= name %>"><span class="input-group-btn"><button class="btn btn-default" type="button">Save</button></div>'),
    events: {
      'click .delete': 'deleteItem',
      'click .edit': 'editItem',
      'click button': 'saveItem',
    },
    initialize: function() {
      this.listenTo(this.model, 'remove', this.remove);
      this.listenTo(this.model, 'change', this.render);
    },
    editItem: function() {
      // Render the edit template
      this.$el.html(this.editTemplate(this.model.toJSON()));
    },
    deleteItem: function() {
      // Delete the item, but wait for the server to respond
      this.model.destroy({wait: true});
    },
    saveItem: function() {
      var name = $.trim(this.$('input').val());
      if (!name) {
        $('#errors').prepend(new Error({message: 'Empty items cannot be saved'}).el);
        // Re-render the original template
        this.$el.html(this.template(this.model.toJSON()));
        return;
      }

      // If name equals the old name, just re-render
      if (name === this.model.get('name')) {
        this.$el.html(this.template(this.model.toJSON()));
        return;
      }

      this.model.set('name', name);
      this.model.save();

      // TODO re-render the original model only on success
    },
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  });

  var User = Backbone.Model.extend({urlRoot: 'users'});
  var Users = Backbone.Collection.extend({
    model: User,
    url: 'users',
    initialize: function() {},
  });

  var Thing = Backbone.Model.extend({urlRoot: 'things'});
  var Things = Backbone.Collection.extend({
    model: Thing,
    url: 'things',
    comparator: function(m) {
      return m.get('timestamp');
    }
  });

  var UserList = Backbone.View.extend({
    el: '#users',
    initialize: function() {
      this.listenTo(this.collection, 'reset add remove', this.render);
      // Render the initial state
      this.render();
    },
    render: function() {
      this.$el.empty();
      _.each(this.collection.models, function(user) {
        // Assign colors as a mod of the user id
        var color = colors[user.get('id') % colors.length];
        // TODO Use a template
        this.$el.append('<li><div class="user" style="background-color:' + color + '"></div></li>');
      }, this);

      // Add the user count
      var len = this.collection.length;
      var userLabel = (len > 1) ? (String(len) + ' Users') : '1 User';
      this.$el.append('<li>' + userLabel + '</li>');
      return this;
    },
  });

  return module;
}();

// Initialize on DOM ready
$(ListOfThings.onready);
