/**
 * 
 */

joint.shapes.basic.Generic.define('crl.BaseElement', {
    attrs: {
        rect: { 'width': 200 },

        '.uml-class-name-rect': { 'stroke': 'black', 'stroke-width': 2, 'fill': '#ffffff' },
        '.uml-class-attrs-rect': { 'stroke': 'black', 'stroke-width': 2, 'fill': '#ffffff' },
        '.uml-class-methods-rect': { 'stroke': 'black', 'stroke-width': 2, 'fill': '#ffffff' },

        '.uml-class-name-text': {
            'ref': '.uml-class-name-rect',
            'ref-y': .5,
            'ref-x': .5,
            'text-anchor': 'middle',
            'y-alignment': 'middle',
            'font-weight': 'bold',
            'fill': 'black',
            'font-size': 12,
            'font-family': 'Times New Roman'
        },
        '.uml-class-attrs-text': {
            'ref': '.uml-class-attrs-rect', 'ref-y': 5, 'ref-x': 5,
            'fill': 'black', 'font-size': 12, 'font-family': 'Times New Roman'
        },
        '.uml-class-methods-text': {
            'ref': '.uml-class-methods-rect', 'ref-y': 5, 'ref-x': 5,
            'fill': 'black', 'font-size': 12, 'font-family': 'Times New Roman'
        }
    },

    name: [],
    attributes: [],
    methods: []
}, {
    markup: [
        '<g class="rotatable">',
        '<g class="scalable">',
        '<rect class="uml-class-name-rect"/><rect class="uml-class-attrs-rect"/><rect class="uml-class-methods-rect"/>',
        '</g>',
        '<text class="uml-class-name-text"/><text class="uml-class-attrs-text"/><text class="uml-class-methods-text"/>',
        '</g>'
    ].join(''),

    initialize: function() {

        this.on('change:name change:attributes change:methods', function() {
            this.updateRectangles();
            this.trigger('uml-update');
        }, this);

        this.updateRectangles();

        joint.shapes.basic.Generic.prototype.initialize.apply(this, arguments);
    },

    getClassName: function() {
        return this.get('name');
    },

    updateRectangles: function() {

        var attrs = this.get('attrs');

        var rects = [
            { type: 'name', text: this.getClassName() },
            { type: 'attrs', text: this.get('attributes') },
            { type: 'methods', text: this.get('methods') }
        ];

        var offsetY = 0;

        rects.forEach(function(rect) {

            var lines = Array.isArray(rect.text) ? rect.text : [rect.text];
            var rectHeight = lines.length * 20 + 20;

            attrs['.uml-class-' + rect.type + '-text'].text = lines.join('\n');
            attrs['.uml-class-' + rect.type + '-rect'].height = rectHeight;
            attrs['.uml-class-' + rect.type + '-rect'].transform = 'translate(0,' + offsetY + ')';

            offsetY += rectHeight;
        });
    }

});
