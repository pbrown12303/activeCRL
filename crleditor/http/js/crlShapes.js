/**
 * 
 */

joint.shapes.basic.Generic.define('crl.Element',
    {
        attrs: {
            rect: {
                width: 300
            },
            '.bounding-rect': {
                stroke: "black",
                'stroke-width': 2,
                fill: "#ffffff",
                height: 40,
                width: 200,
                transform: "translate(0,0)"
            },
            '.image': {
                'type':'image',
                'ref-x': 1.0,
                'ref-y': 1.0,
                ref: ".bounding-rect",
                width: 16,
                height: 16,
                'xlink:href': ""
            },
            '.abstractions-text': {
                ref: ".bounding-rect",
                'ref-y': 1.0,
                'ref-x': 1.0 + 18,
                'text-anchor': "left",
                'y-alignment': "top",
                'font-weight': "normal",
                'font-style': "italic",
                fill: "black",
                'font-size': 10,
                'font-family': "Go,  Helvetica, Ariel, sans-serif",
                'text':"defaultAbstractionsText"
            },
            '.labelText': {
                ref: ".bounding-rect",
                'ref-y': 4.0 + 16,
                'ref-x': 3.0,
                'text-anchor': "left",
                'y-alignment': "top",
                'font-weight': "bold",
                fill: "black",
                'font-size': 12,
                'font-family': "Go,  Helvetica, Ariel, sans-serif",
                'text':""
            }
        },
        crlJointID: "",
        name:  "labelDefault", 
        abstractions: "defaultAbstractions",
        icon:"",
        displayLabelYOffset:0.0
    },
    {
        markup: "<g class=\"rotatable\">" +
        "<g class=\"scalable\">" +
        "<rect class=\"bounding-rect\"/>" +
        "</g>" +
        "<image class=\"image\"/>" +
        "<text class=\"abstractions-text\"/>" + 
        "<text class=\"labelText\" />" +
        "</g>",

        initialize: function () {

            this.on('change:name change:abstractions change:icon change:displayLabelYOffset', function () {
               this.updateRectangles();
               this.trigger('crl-update');
            }, this);

            this.updateRectangles();

            joint.shapes.basic.Generic.prototype.initialize.apply(this, arguments);
        },

        // getClassName: function () {
        //     return this.get('name');
        // },

        updateRectangles: function () {
            var attrs = this.get('attrs');
            var boundingRectAttr = attrs['.bounding-rect'];
            attrs['.labelText'].text = this.get("name");
            attrs['.image']['xlink:href'] = this.get("icon");
            attrs['.abstractions-text'].text = this.get("abstractions");
            attrs['.labelText']['ref-y'] = this.get("displayLabelYOffset");
            this.resize(boundingRectAttr.width , boundingRectAttr.height);


            // var attrs = this.get('attrs');

            // var rects = [
            //     { type: 'name', text: this.getClassName() } /*,
            // { type: 'attrs', text: this.get('attributes') },
            // { type: 'methods', text: this.get('methods') } */
            // ];

            // var offsetY = 0;

            // var lines = [this.getClassName()];
            // var rectHeight = 1 * 12 + 6;

            // attrs['.label-text'].text = lines.join('\n');
            // attrs['.label-rect'].height = rectHeight;
            // var rectWidth = calculateTextWidth(attrs['.label-text'].text) + 6;
            // attrs['.label-rect'].transform = 'translate(0,' + offsetY + ')';
            // this.resize(rectWidth, rectHeight);

            // offsetY += rectHeight;

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

    joint.shapes.crl.ElementView = joint.dia.ElementView.extend({}, {

        initialize: function() {
    
            joint.dia.ElementView.prototype.initialize.apply(this, arguments);
    
            this.listenTo(this.model, 'crl-update', function() {
                this.update();
                this.resize();
            });
        }
    });
    