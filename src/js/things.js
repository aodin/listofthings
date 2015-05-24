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

  var WEBSOCKET_URI = 'ws://' + document.URL.split('/', 3)[2] + '/feeds/v1/things';

  var App = Backbone.View.extend({
    el: '#main',
    initialize: function(options) {
      // Attach the user and thing collections
      this.users = options.users;
      this.things = options.things;

      // Create views that listen
      // TODO pass the errors handler to each list
      new UserList({collection: this.users});
      new ThingsList({collection: this.things});

      // Cache DOM elements
      this.$errors = $('errors');

      // Create a new websocket
      this.ws = new WebSocket(WEBSOCKET_URI);
      // this.ws.onopen = this.join.bind(this);
      this.ws.onmessage = this.onMessage.bind(this);
      this.ws.onerror = this.onError.bind(this);
      this.ws.onclose = this.leave.bind(this);
    },
    onMessage: function(msg) {
      // Translate the message as JSON
      var payload = JSON.parse(msg.data);

      // TODO common/whitelist store of resources
      this.handleEvent(this[payload.resource], payload.method, payload.content);
    },
    handleEvent: function(collection, method, content) {
      console.log('handling:', collection, method, content); 
      switch (method) {
        case 'LIST':
          collection.reset(content);
          break;
        case 'CREATE':
          collection.add(content);
          break;
        case 'DELETE':
          collection.remove(content);
          break;
        case 'UPDATE':
          collection.add(content, {merge: true});
          break;
      }
    },
    onError: function() {},
    leave: function() {},
    sync: function(method, model) {
      // Check ready state
      if (this.ws.readyState !== 1) {
        // Return after displaying an error
        this.$error.prepend(new Error({message: 'Could not connect to server'}).el);
        return;
      }

      // TODO translate the messages here?
      console.log('sending:', model.collection.url, method);
      var msg = {
        resource: model.collection.url,
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
    var app = new App({things: new Things(), users: new Users()});

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

  var ThingsList = Backbone.View.extend({
    el: '#things',
    events: {
      'keyup #create-name': 'proxyEnter',
      'click #create': 'createItem',
    },
    initialize: function() {
      this.listenTo(this.collection, 'reset', this.render);
      this.listenTo(this.collection, 'add', this.renderItem);
    },
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

      // Create the new thing
      this.collection.create({name: $input.val()}, {wait: true});

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
    template: _.template('<h3><%- name %> <span class="edit"><small>edit</small></span><span class="delete"><small>delete</small></span></h3>'),
    editTemplate: _.template('<div class="input-group"><input type="text" class="form-control" value="<%- name %>"><span class="input-group-btn"><button class="btn btn-default" type="button">Save</button></div>'),
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
