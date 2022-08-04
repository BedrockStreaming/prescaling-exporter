// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"

	scheme "github.com/bedrockstreaming/prescaling-exporter/generated/client/clientset/versioned/scheme"
	v1 "github.com/bedrockstreaming/prescaling-exporter/pkg/apis/prescaling.bedrock.tech/v1"
)

// PrescalingEventsGetter has a method to return a PrescalingEventInterface.
// A group's client should implement this interface.
type PrescalingEventsGetter interface {
	PrescalingEvents(namespace string) PrescalingEventInterface
}

// PrescalingEventInterface has methods to work with PrescalingEvent resources.
type PrescalingEventInterface interface {
	Create(ctx context.Context, prescalingEvent *v1.PrescalingEvent, opts metav1.CreateOptions) (*v1.PrescalingEvent, error)
	Update(ctx context.Context, prescalingEvent *v1.PrescalingEvent, opts metav1.UpdateOptions) (*v1.PrescalingEvent, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.PrescalingEvent, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PrescalingEventList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.PrescalingEvent, err error)
	PrescalingEventExpansion
}

// prescalingEvents implements PrescalingEventInterface
type prescalingEvents struct {
	client rest.Interface
	ns     string
}

// newPrescalingEvents returns a PrescalingEvents
func newPrescalingEvents(c *PrescalingV1Client, namespace string) *prescalingEvents {
	return &prescalingEvents{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the prescalingEvent, and returns the corresponding prescalingEvent object, and an error if there is any.
func (c *prescalingEvents) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.PrescalingEvent, err error) {
	result = &v1.PrescalingEvent{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("prescalingevents").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PrescalingEvents that match those selectors.
func (c *prescalingEvents) List(ctx context.Context, opts metav1.ListOptions) (result *v1.PrescalingEventList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.PrescalingEventList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("prescalingevents").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested prescalingEvents.
func (c *prescalingEvents) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("prescalingevents").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a prescalingEvent and creates it.  Returns the server's representation of the prescalingEvent, and an error, if there is any.
func (c *prescalingEvents) Create(ctx context.Context, prescalingEvent *v1.PrescalingEvent, opts metav1.CreateOptions) (result *v1.PrescalingEvent, err error) {
	result = &v1.PrescalingEvent{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("prescalingevents").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(prescalingEvent).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a prescalingEvent and updates it. Returns the server's representation of the prescalingEvent, and an error, if there is any.
func (c *prescalingEvents) Update(ctx context.Context, prescalingEvent *v1.PrescalingEvent, opts metav1.UpdateOptions) (result *v1.PrescalingEvent, err error) {
	result = &v1.PrescalingEvent{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("prescalingevents").
		Name(prescalingEvent.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(prescalingEvent).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the prescalingEvent and deletes it. Returns an error if one occurs.
func (c *prescalingEvents) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("prescalingevents").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *prescalingEvents) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("prescalingevents").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched prescalingEvent.
func (c *prescalingEvents) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.PrescalingEvent, err error) {
	result = &v1.PrescalingEvent{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("prescalingevents").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}