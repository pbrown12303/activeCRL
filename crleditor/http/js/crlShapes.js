/**
 * 
 */

joint.shapes.basic.Generic.define('crl.Element', {
    attrs: {
        rect: { width: 300  },

        '.label-rect': { 'stroke': 'black', 'stroke-width': 2, 'fill': '#ffffff' },
        /*        '.uml-class-attrs-rect': { 'stroke': 'black', 'stroke-width': 2, 'fill': '#ffffff' },
                '.uml-class-methods-rect': { 'stroke': 'black', 'stroke-width': 2, 'fill': '#ffffff' },
*/
        '.label-text': {
            'ref': '.label-rect',
            'ref-y': .5,
            'ref-x': .5,
            'text-anchor': 'middle',
            'y-alignment': 'middle',
            'font-weight': 'bold',
            'fill': 'black',
            'font-size': 12,
            'font-family': 'Go,  Helvetica, Ariel, sans-serif'
        },
        '.uml-class-attrs-text': {
            'ref': '.uml-class-attrs-rect', 'ref-y': 5, 'ref-x': 5,
            'fill': 'black', 'font-size': 12, 'font-family': 'Go,  Helvetica, Ariel, sans-serif' 
        },
        '.uml-class-methods-text': {
            'ref': '.uml-class-methods-rect', 'ref-y': 5, 'ref-x': 5,
            'fill': 'black', 'font-size': 12, 'font-family': 'Go,  Helvetica, Ariel, sans-serif'
        }
    },

    name: [],
    attributes: [],
    methods: []
}, {
    markup: [
        '<g class="rotatable">',
        '<g class="scalable">',
        '<rect class="label-rect"/><rect class="uml-class-attrs-rect"/><rect class="uml-class-methods-rect"/>',
        '</g>',
        '<text class="label-text"/><text class="uml-class-attrs-text"/><text class="uml-class-methods-text"/>',
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
            { type: 'name', text: this.getClassName() } /*,
            { type: 'attrs', text: this.get('attributes') },
            { type: 'methods', text: this.get('methods') } */
        ];

        var offsetY = 0;
        
        var lines = [this.getClassName()]; 
        var rectHeight = 1 * 12 + 6; 

        attrs['.label-text'].text = lines.join('\n');
        attrs['.label-rect'].height = rectHeight;
        var rectWidth = calculateTextWidth(attrs['.label-text'].text) + 6;
        attrs['.label-rect'].transform = 'translate(0,' + offsetY + ')';
        this.resize(rectWidth, rectHeight);

        offsetY += rectHeight;

/*        rects.forEach(function(rect) {

            var lines = Array.isArray(rect.text) ? rect.text : [rect.text];
            var rectHeight = lines.length * 20 + 20;

            attrs['.uml-class-' + rect.type + '-text'].text = lines.join('\n');
            attrs['.uml-class-' + rect.type + '-rect'].height = rectHeight;
            attrs['.uml-class-' + rect.type + '-rect'].transform = 'translate(0,' + offsetY + ')';

            offsetY += rectHeight;
        });
        */
    }

});
