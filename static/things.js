var ws_uri = 'ws://' + document.URL.split('/', 3)[2] + '/events';

var Error = Backbone.View.extend({
  tagName: 'li',
  initialize: function(options) {
    this.message = options.message;
    // Destory self after a timeout
    var self = this;
    setTimeout(function() {self.remove()}, 4000);
    return this.render();
  },
  render: function() {
    console.log('rendering new error:', this.message);
    this.$el.html(this.message);
    return this
  },
});

var List = Backbone.View.extend({
  el: '#things',
  events: {
    'keyup #create-name': 'proxyEnter',
    'click #create': 'createItem',
  },
  initialize: function(options) {
    _.bindAll(this, 'onmessage', 'onopen');
    this.listenTo(this.collection, 'reset', this.render)
    this.listenTo(this.collection, 'add', this.renderItem)
  },
  onmessage: function(result) {
    console.log('Message:', result);
    var data = JSON.parse(result.data);
    if (!data) return;

    // TODO a switch makes sense here
    if (data.body === 'init') {
      // Initial bootstrap message
      console.log('content:', data.content);
      // TODO This should be behavior of Backbone.sync
      this.collection.reset(data.content)
    } else if (data.body === 'create') {
      // Item creation message
      this.collection.add(data.content)
    }  else if (data.body === 'delete') {
      // Item deletion message
      // Get the item from the collection and call delete
      // TODO This should be behavior of Backbone.sync
      var thing = this.collection.get(data.content.id);
      console.log('thing to delete:', thing);
      this.collection.remove(thing);
    }  else if (data.body === 'update') {
      // Item update message
      // TODO This should be behavior of Backbone.sync
      this.collection.add(data.content, {merge: true});

    } else if (data.body === 'error') {
      // Error message, display for a timeout
      this.$('#errors').prepend(new Error({message: data.content}).el);
    }
  },
  onopen: function(data) {
    console.log('Socket open');
  },
  proxyEnter: function(e) {
    if (e.keyCode == 13) this.createItem();
  },
  createItem: function() {
    var $input = this.$('#create-name');

    // Trim whitespace
    var name = $.trim($input.val());
    if (!name) {
      this.$('#errors').prepend(new Error({message: "Empty items cannot be created"}).el);
      return
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
    return 
  },
  renderItem: function(m) {
    var $list = this.$('ol');
    var item = new Item({model: m});
    $list.append(item.render().el);
  }
});

var Item = Backbone.View.extend({
  tagName: 'li',
  template: '<h3><%= name %> <span class="edit"><small>edit</small></span><span class="delete"><small>delete</small></span></h3>',
  editTemplate: '<div class="input-group"><input type="text" class="form-control" value="<%= name %>"><span class="input-group-btn"><button class="btn btn-default" type="button">Save</button></div>',
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
    this.$el.html(_.template(this.editTemplate, this.model.toJSON()));
  },
  deleteItem: function() {
    // Delete the item, but wait for the server to respond
    console.log('destroying model');
    this.model.destroy({wait: true});
  },
  saveItem: function() {
    console.log('saving item');
    var name = $.trim(this.$('input').val());
    if (!name) {
      $('#errors').prepend(new Error({message: "Empty items cannot be saved"}).el);
      // Re-render the original template
      this.$el.html(_.template(this.template, this.model.toJSON()));
      return;
    }

    // If name equals the old name, just re-render
    if (name === this.model.get('name')) {
      this.$el.html(_.template(this.template, this.model.toJSON()));
      return;
    }

    console.log('setting name:', name);
    this.model.set('name', name);
    this.model.save();

    // TODO re-render the original model only on success
  },
  render: function() {
    this.$el.html(_.template(this.template, this.model.toJSON()));
    return this
  }
});

var Thing = Backbone.Model.extend({});
var Things = Backbone.Collection.extend({
  model: Thing,
});

function CreateWebsocketSync(ws) {
  var WebsocketSync = function(method, model, options) {
    console.log('sync:', method, model, options);
    console.log(model.toJSON());
    ws.send(JSON.stringify({method: method, content: model.toJSON()}));
  }
  return WebsocketSync;
}

$(function() {
  var things = new Things();
  var list = new List({collection: things});

  console.log("Connecting to:", ws_uri);
  var w = new WebSocket(ws_uri);

  // Overwrite the sync method
  Backbone.sync = CreateWebsocketSync(w);

  // TODO Better event handler
  w.onopen = list.onopen;
  w.onmessage = list.onmessage;
});